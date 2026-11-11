package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/ihulsbus/cookbook/shared/oidc"

	m "user-service/internal/models"
)

// userServiceMock is a UserService whose behavior each test configures
// via its function fields, rather than switching on magic input values.
type userServiceMock struct {
	searchUser       func(ctx context.Context, query string) (m.UserDTO, error)
	getDisplayName   func(ctx context.Context, id uuid.UUID) (m.UserDTO, error)
	listMyHouseholds func(ctx context.Context, userID uuid.UUID) ([]m.HouseholdDTO, error)
	createHousehold  func(ctx context.Context, ownerID uuid.UUID, name string) (m.HouseholdDTO, error)
	addMember        func(ctx context.Context, householdID string, requesterID uuid.UUID, query string) (m.HouseholdDTO, error)
	removeMember     func(ctx context.Context, householdID string, requesterID, targetUserID uuid.UUID) error
	leaveHousehold   func(ctx context.Context, householdID string, userID uuid.UUID) error
	deleteHousehold  func(ctx context.Context, householdID string, requesterID uuid.UUID) error
}

func (s *userServiceMock) SearchUser(ctx context.Context, query string) (m.UserDTO, error) {
	return s.searchUser(ctx, query)
}
func (s *userServiceMock) GetDisplayName(ctx context.Context, id uuid.UUID) (m.UserDTO, error) {
	return s.getDisplayName(ctx, id)
}
func (s *userServiceMock) ListMyHouseholds(ctx context.Context, userID uuid.UUID) ([]m.HouseholdDTO, error) {
	return s.listMyHouseholds(ctx, userID)
}
func (s *userServiceMock) CreateHousehold(ctx context.Context, ownerID uuid.UUID, name string) (m.HouseholdDTO, error) {
	return s.createHousehold(ctx, ownerID, name)
}
func (s *userServiceMock) AddMember(ctx context.Context, householdID string, requesterID uuid.UUID, query string) (m.HouseholdDTO, error) {
	return s.addMember(ctx, householdID, requesterID, query)
}
func (s *userServiceMock) RemoveMember(ctx context.Context, householdID string, requesterID, targetUserID uuid.UUID) error {
	return s.removeMember(ctx, householdID, requesterID, targetUserID)
}
func (s *userServiceMock) LeaveHousehold(ctx context.Context, householdID string, userID uuid.UUID) error {
	return s.leaveHousehold(ctx, householdID, userID)
}
func (s *userServiceMock) DeleteHousehold(ctx context.Context, householdID string, requesterID uuid.UUID) error {
	return s.deleteHousehold(ctx, householdID, requesterID)
}

var testRequesterID = uuid.New()

// authedContext builds a gin test context with testRequesterID set as
// the authenticated principal, as oidc.Module.RequireAuth would before
// any handler in this package runs.
func authedContext(w *httptest.ResponseRecorder, req *http.Request) *gin.Context {
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	oidc.SetPrincipal(c, &oidc.Principal{UserID: testRequesterID.String()})
	return c
}

func newHandlers(mock *userServiceMock) *HttpHandlers {
	return NewHttpHandlers(mock, log.New())
}

func TestSearchUsers_OK(t *testing.T) {
	mock := &userServiceMock{
		searchUser: func(ctx context.Context, query string) (m.UserDTO, error) {
			assert.Equal(t, "alice", query)
			return m.UserDTO{ID: uuid.New(), Username: "alice", DisplayName: "Alice"}, nil
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/user/users?query=alice", nil)
	c := authedContext(w, req)

	h.SearchUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchUsers_NotFound(t *testing.T) {
	mock := &userServiceMock{
		searchUser: func(ctx context.Context, query string) (m.UserDTO, error) {
			return m.UserDTO{}, errors.New("user not found")
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/user/users?query=nobody", nil)
	c := authedContext(w, req)

	h.SearchUsers(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUser_InvalidID(t *testing.T) {
	h := newHandlers(&userServiceMock{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/user/users/not-a-uuid", nil)
	c := authedContext(w, req)
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}

	h.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUser_OK(t *testing.T) {
	id := uuid.New()
	mock := &userServiceMock{
		getDisplayName: func(ctx context.Context, got uuid.UUID) (m.UserDTO, error) {
			assert.Equal(t, id, got)
			return m.UserDTO{ID: id, DisplayName: "Alice"}, nil
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/user/users/"+id.String(), nil)
	c := authedContext(w, req)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	h.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListMyHouseholds_OK(t *testing.T) {
	mock := &userServiceMock{
		listMyHouseholds: func(ctx context.Context, userID uuid.UUID) ([]m.HouseholdDTO, error) {
			assert.Equal(t, testRequesterID, userID)
			return []m.HouseholdDTO{{ID: "g-1", Name: "The Smiths"}}, nil
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v2/user/me/households", nil)
	c := authedContext(w, req)

	h.ListMyHouseholds(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateHousehold_OK(t *testing.T) {
	mock := &userServiceMock{
		createHousehold: func(ctx context.Context, ownerID uuid.UUID, name string) (m.HouseholdDTO, error) {
			assert.Equal(t, testRequesterID, ownerID)
			assert.Equal(t, "The Smiths", name)
			return m.HouseholdDTO{ID: "g-1", Name: name, OwnerID: ownerID}, nil
		},
	}
	h := newHandlers(mock)

	body, _ := json.Marshal(createHouseholdRequest{Name: "The Smiths"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/user/households", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c := authedContext(w, req)

	h.CreateHousehold(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestCreateHousehold_NameRequired(t *testing.T) {
	mock := &userServiceMock{
		createHousehold: func(ctx context.Context, ownerID uuid.UUID, name string) (m.HouseholdDTO, error) {
			return m.HouseholdDTO{}, errors.New("name is required")
		},
	}
	h := newHandlers(mock)

	body, _ := json.Marshal(createHouseholdRequest{Name: ""})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/user/households", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c := authedContext(w, req)

	h.CreateHousehold(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddMember_Forbidden(t *testing.T) {
	householdID := uuid.New()
	mock := &userServiceMock{
		addMember: func(ctx context.Context, householdID string, requesterID uuid.UUID, query string) (m.HouseholdDTO, error) {
			return m.HouseholdDTO{}, errors.New("forbidden")
		},
	}
	h := newHandlers(mock)

	body, _ := json.Marshal(addMemberRequest{Query: "bob"})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/user/households/"+householdID.String()+"/members", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c := authedContext(w, req)
	c.Params = gin.Params{{Key: "id", Value: householdID.String()}}

	h.AddMember(c)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRemoveMember_OK(t *testing.T) {
	householdID := uuid.New()
	targetID := uuid.New()
	mock := &userServiceMock{
		removeMember: func(ctx context.Context, gotHouseholdID string, requesterID, targetUserID uuid.UUID) error {
			assert.Equal(t, householdID.String(), gotHouseholdID)
			assert.Equal(t, testRequesterID, requesterID)
			assert.Equal(t, targetID, targetUserID)
			return nil
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/user/households/"+householdID.String()+"/members/"+targetID.String(), nil)
	c := authedContext(w, req)
	c.Params = gin.Params{{Key: "id", Value: householdID.String()}, {Key: "userId", Value: targetID.String()}}

	h.RemoveMember(c)

	// ctx.Status() doesn't flush to the ResponseRecorder until the first
	// body write (which a 204 never has), so assert on the Writer's
	// status instead of w.Code.
	assert.Equal(t, http.StatusNoContent, c.Writer.Status())
}

func TestLeaveHousehold_OwnerCannotLeave(t *testing.T) {
	householdID := uuid.New()
	mock := &userServiceMock{
		leaveHousehold: func(ctx context.Context, householdID string, userID uuid.UUID) error {
			return errors.New("owner cannot leave")
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/user/households/"+householdID.String()+"/leave", nil)
	c := authedContext(w, req)
	c.Params = gin.Params{{Key: "id", Value: householdID.String()}}

	h.LeaveHousehold(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteHousehold_OK(t *testing.T) {
	householdID := uuid.New()
	mock := &userServiceMock{
		deleteHousehold: func(ctx context.Context, gotHouseholdID string, requesterID uuid.UUID) error {
			assert.Equal(t, householdID.String(), gotHouseholdID)
			return nil
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/user/households/"+householdID.String(), nil)
	c := authedContext(w, req)
	c.Params = gin.Params{{Key: "id", Value: householdID.String()}}

	h.DeleteHousehold(c)

	assert.Equal(t, http.StatusNoContent, c.Writer.Status())
}

func TestDeleteHousehold_NotFound(t *testing.T) {
	householdID := uuid.New()
	mock := &userServiceMock{
		deleteHousehold: func(ctx context.Context, gotHouseholdID string, requesterID uuid.UUID) error {
			return errors.New("not found")
		},
	}
	h := newHandlers(mock)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v2/user/households/"+householdID.String(), nil)
	c := authedContext(w, req)
	c.Params = gin.Params{{Key: "id", Value: householdID.String()}}

	h.DeleteHousehold(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
