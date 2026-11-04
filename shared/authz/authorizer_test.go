package authz

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ihulsbus/cookbook/shared/oidc"
)

type loggerMock struct{}

func (loggerMock) Warnf(format string, args ...interface{})  {}
func (loggerMock) Errorf(format string, args ...interface{}) {}

func sharesJSON(t *testing.T, shares []rmq.ShareEntry) []byte {
	t.Helper()
	b, err := json.Marshal(shares)
	if err != nil {
		t.Fatalf("failed to marshal shares: %v", err)
	}
	return b
}

// TestEvaluate_RuleMatrix is the rule matrix KIT-27 explicitly asks for:
// owner read/write, share-granted read/write (user and group), public
// read, and deny-by-default for everyone else.
func TestEvaluate_RuleMatrix(t *testing.T) {
	owner := "owner-1"
	stranger := "stranger-1"
	sharedUser := "shared-user-1"
	groupMember := "group-member-1"

	tests := []struct {
		name    string
		row     *RecipeAccess
		userID  string
		groups  []string
		perm    Permission
		allowed bool
	}{
		{"owner can read", &RecipeAccess{OwnerID: owner, Visibility: "private"}, owner, nil, PermissionRead, true},
		{"owner can write", &RecipeAccess{OwnerID: owner, Visibility: "private"}, owner, nil, PermissionWrite, true},
		{"stranger denied read on private", &RecipeAccess{OwnerID: owner, Visibility: "private"}, stranger, nil, PermissionRead, false},
		{"stranger denied write on private", &RecipeAccess{OwnerID: owner, Visibility: "private"}, stranger, nil, PermissionWrite, false},
		{"stranger can read public", &RecipeAccess{OwnerID: owner, Visibility: "public"}, stranger, nil, PermissionRead, true},
		{"stranger cannot write public", &RecipeAccess{OwnerID: owner, Visibility: "public"}, stranger, nil, PermissionWrite, false},
		{
			"user share with read grants read",
			&RecipeAccess{OwnerID: owner, Visibility: "shared", Shares: sharesJSON(t, []rmq.ShareEntry{
				{TargetType: "user", TargetID: sharedUser, Permission: "read"},
			})},
			sharedUser, nil, PermissionRead, true,
		},
		{
			"user share with read does not grant write",
			&RecipeAccess{OwnerID: owner, Visibility: "shared", Shares: sharesJSON(t, []rmq.ShareEntry{
				{TargetType: "user", TargetID: sharedUser, Permission: "read"},
			})},
			sharedUser, nil, PermissionWrite, false,
		},
		{
			"user share with write grants write",
			&RecipeAccess{OwnerID: owner, Visibility: "shared", Shares: sharesJSON(t, []rmq.ShareEntry{
				{TargetType: "user", TargetID: sharedUser, Permission: "write"},
			})},
			sharedUser, nil, PermissionWrite, true,
		},
		{
			"group share grants read via membership",
			&RecipeAccess{OwnerID: owner, Visibility: "shared", Shares: sharesJSON(t, []rmq.ShareEntry{
				{TargetType: "group", TargetID: "household", Permission: "read"},
			})},
			groupMember, []string{"household"}, PermissionRead, true,
		},
		{
			"group share grants nothing without membership",
			&RecipeAccess{OwnerID: owner, Visibility: "shared", Shares: sharesJSON(t, []rmq.ShareEntry{
				{TargetType: "group", TargetID: "household", Permission: "read"},
			})},
			groupMember, []string{"other-group"}, PermissionRead, false,
		},
		{
			"unrelated share does not grant access",
			&RecipeAccess{OwnerID: owner, Visibility: "shared", Shares: sharesJSON(t, []rmq.ShareEntry{
				{TargetType: "user", TargetID: "someone-else", Permission: "write"},
			})},
			stranger, nil, PermissionRead, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			principal := &oidc.Principal{UserID: tt.userID, Groups: tt.groups}

			allowed, err := evaluate(tt.row, principal, tt.perm)

			assert.NoError(t, err)
			assert.Equal(t, tt.allowed, allowed)
		})
	}
}

func TestAuthorizer_CanRead_ProjectionMiss_NoFallback(t *testing.T) {
	db, mock := newMockDatabase(t)
	a := NewAuthorizer(db, loggerMock{})

	recipeID := uuid.New()
	principal := &oidc.Principal{UserID: "owner-1"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnError(gorm.ErrRecordNotFound)

	allowed, err := a.CanRead(context.Background(), principal, recipeID)

	assert.False(t, allowed)
	assert.ErrorIs(t, err, ErrProjectionMiss)
	assert.Equal(t, uint64(1), a.MissCount())
	assert.Equal(t, uint64(0), a.HitCount())
}

// TestAuthorizer_CanRead_ProjectionMiss_FallbackResolves is KIT-28's
// contract from shared/authz's side: when the projection is missing,
// SetMissFallback's function is called, and if it succeeds the lookup is
// retried once rather than failing outright.
func TestAuthorizer_CanRead_ProjectionMiss_FallbackResolves(t *testing.T) {
	db, mock := newMockDatabase(t)
	a := NewAuthorizer(db, loggerMock{})

	recipeID := uuid.New()
	principal := &oidc.Principal{UserID: "owner-1"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnError(gorm.ErrRecordNotFound)

	fallbackCalled := false
	a.SetMissFallback(func(ctx context.Context, id uuid.UUID) error {
		fallbackCalled = true
		assert.Equal(t, recipeID, id)
		return nil
	})

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "owner_id", "visibility", "version", "deleted"}).
			AddRow(recipeID, "owner-1", "private", 1000, false))

	allowed, err := a.CanRead(context.Background(), principal, recipeID)

	assert.True(t, fallbackCalled)
	assert.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, uint64(1), a.MissCount())
}

func TestAuthorizer_CanRead_ProjectionMiss_FallbackFails(t *testing.T) {
	db, mock := newMockDatabase(t)
	a := NewAuthorizer(db, loggerMock{})

	recipeID := uuid.New()
	principal := &oidc.Principal{UserID: "owner-1"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnError(gorm.ErrRecordNotFound)

	a.SetMissFallback(func(ctx context.Context, id uuid.UUID) error {
		return errors.New("recipe-service unreachable")
	})

	allowed, err := a.CanRead(context.Background(), principal, recipeID)

	assert.False(t, allowed)
	assert.ErrorIs(t, err, ErrProjectionMiss)
}

func TestAuthorizer_CanRead_DeletedRecipe(t *testing.T) {
	db, mock := newMockDatabase(t)
	a := NewAuthorizer(db, loggerMock{})

	recipeID := uuid.New()
	principal := &oidc.Principal{UserID: "owner-1"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "recipe_accesses" WHERE recipe_id = $1`)).
		WillReturnRows(sqlmock.NewRows([]string{"recipe_id", "owner_id", "visibility", "version", "deleted"}).
			AddRow(recipeID, "owner-1", "private", 1000, true))

	allowed, err := a.CanRead(context.Background(), principal, recipeID)

	assert.False(t, allowed)
	assert.Error(t, err)
	assert.NotErrorIs(t, err, ErrProjectionMiss)
}
