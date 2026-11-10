package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"user-service/internal/authentik"
)

// fakeAuthentik is an in-memory AuthentikClient for tests - no network,
// no real Authentik instance required.
type fakeAuthentik struct {
	usersByUUID map[string]authentik.User
	usersByName map[string]authentik.User
	usersByPK   map[int]authentik.User
	groups      map[string]authentik.Group
	nextGroupPK int
}

func newFakeAuthentik() *fakeAuthentik {
	return &fakeAuthentik{
		usersByUUID: map[string]authentik.User{},
		usersByName: map[string]authentik.User{},
		usersByPK:   map[int]authentik.User{},
		groups:      map[string]authentik.Group{},
	}
}

func (f *fakeAuthentik) addUser(u authentik.User) {
	f.usersByUUID[u.UUID] = u
	f.usersByPK[u.PK] = u
	f.usersByName[u.Username] = u
	f.usersByName[u.Email] = u
}

func (f *fakeAuthentik) FindUserByUUID(_ context.Context, userUUID string) (authentik.User, error) {
	u, ok := f.usersByUUID[userUUID]
	if !ok {
		return authentik.User{}, authentik.ErrNotFound
	}
	return u, nil
}

func (f *fakeAuthentik) FindUserExact(_ context.Context, query string) (authentik.User, error) {
	u, ok := f.usersByName[query]
	if !ok {
		return authentik.User{}, authentik.ErrNotFound
	}
	return u, nil
}

func (f *fakeAuthentik) GetUserByPK(_ context.Context, pk int) (authentik.User, error) {
	u, ok := f.usersByPK[pk]
	if !ok {
		return authentik.User{}, authentik.ErrNotFound
	}
	return u, nil
}

func (f *fakeAuthentik) FindGroupByName(_ context.Context, name, parentPK string) (authentik.Group, error) {
	for _, g := range f.groups {
		if g.Name == name && g.Parent == parentPK {
			return g, nil
		}
	}
	return authentik.Group{}, authentik.ErrNotFound
}

func (f *fakeAuthentik) GetGroup(_ context.Context, pk string) (authentik.Group, error) {
	g, ok := f.groups[pk]
	if !ok {
		return authentik.Group{}, authentik.ErrNotFound
	}
	return g, nil
}

func (f *fakeAuthentik) CreateGroup(_ context.Context, g authentik.Group) (authentik.Group, error) {
	f.nextGroupPK++
	g.PK = uuid.New().String()
	f.groups[g.PK] = g
	return g, nil
}

func (f *fakeAuthentik) UpdateGroupUsers(_ context.Context, pk string, userPKs []int) error {
	g, ok := f.groups[pk]
	if !ok {
		return authentik.ErrNotFound
	}
	g.Users = userPKs
	f.groups[pk] = g
	return nil
}

func (f *fakeAuthentik) DeleteGroup(_ context.Context, pk string) error {
	if _, ok := f.groups[pk]; !ok {
		return authentik.ErrNotFound
	}
	delete(f.groups, pk)
	return nil
}

// fakeRepo is an in-memory UserRepository for tests.
type fakeRepo struct {
	names map[uuid.UUID]string
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{names: map[uuid.UUID]string{}}
}

func (f *fakeRepo) GetDisplayName(id uuid.UUID) (string, error) {
	n, ok := f.names[id]
	if !ok {
		return "", errors.New("user not found")
	}
	return n, nil
}

func (f *fakeRepo) UpsertDisplayName(id uuid.UUID, displayName string) error {
	f.names[id] = displayName
	return nil
}

func newTestService() (*UserService, *fakeAuthentik, *fakeRepo) {
	ak := newFakeAuthentik()
	repo := newFakeRepo()
	logger := log.New()
	return NewUserService(ak, repo, logger), ak, repo
}

func TestSearchUser(t *testing.T) {
	s, ak, repo := newTestService()
	id := uuid.New()
	ak.addUser(authentik.User{PK: 1, UUID: id.String(), Username: "alice", Name: "Alice A", Email: "alice@example.com"})

	dto, err := s.SearchUser(context.Background(), "alice")
	assert.NoError(t, err)
	assert.Equal(t, id, dto.ID)
	assert.Equal(t, "alice", dto.Username)
	assert.Equal(t, "Alice A", dto.DisplayName)
	assert.Equal(t, "Alice A", repo.names[id])
}

func TestSearchUser_QueryRequired(t *testing.T) {
	s, _, _ := newTestService()
	_, err := s.SearchUser(context.Background(), "")
	assert.EqualError(t, err, "query is required")
}

func TestSearchUser_NotFound(t *testing.T) {
	s, _, _ := newTestService()
	_, err := s.SearchUser(context.Background(), "nobody")
	assert.EqualError(t, err, "user not found")
}

func TestGetDisplayName_CacheHit(t *testing.T) {
	s, _, repo := newTestService()
	id := uuid.New()
	repo.names[id] = "Cached Alice"

	dto, err := s.GetDisplayName(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "Cached Alice", dto.DisplayName)
}

func TestGetDisplayName_FallsBackToAuthentik(t *testing.T) {
	s, ak, repo := newTestService()
	id := uuid.New()
	ak.addUser(authentik.User{PK: 1, UUID: id.String(), Username: "bob", Name: "Bob B"})

	dto, err := s.GetDisplayName(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, "Bob B", dto.DisplayName)
	assert.Equal(t, "Bob B", repo.names[id])
}

func TestGetDisplayName_NotFound(t *testing.T) {
	s, _, _ := newTestService()
	_, err := s.GetDisplayName(context.Background(), uuid.New())
	assert.EqualError(t, err, "user not found")
}

// seedHousehold registers a household group with owner as its owner and
// first member, plus a "households" parent group.
func seedHousehold(ak *fakeAuthentik, name string, owner authentik.User, members ...authentik.User) authentik.Group {
	parent, err := ak.FindGroupByName(context.Background(), householdsParentGroupName, "")
	if err != nil {
		parent, _ = ak.CreateGroup(context.Background(), authentik.Group{Name: householdsParentGroupName})
	}

	userPKs := []int{owner.PK}
	for _, m := range members {
		userPKs = append(userPKs, m.PK)
	}

	g, _ := ak.CreateGroup(context.Background(), authentik.Group{
		Name:   name,
		Parent: parent.PK,
		Users:  userPKs,
		Attributes: map[string]any{
			"cookbook_household": true,
			"owner":              owner.UUID,
		},
	})
	return g
}

func TestListMyHouseholds_FiltersNonHouseholdGroups(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	ak.addUser(owner)

	household := seedHousehold(ak, "The Smiths", owner)

	// An unrelated group the user also belongs to (e.g. cookbook-admin)
	// must not show up as a household.
	other, _ := ak.CreateGroup(context.Background(), authentik.Group{Name: "cookbook-admin", Users: []int{owner.PK}})

	ownerWithGroups := owner
	ownerWithGroups.Groups = []string{household.PK, other.PK}
	ak.addUser(ownerWithGroups)

	ownerID, _ := uuid.Parse(owner.UUID)
	households, err := s.ListMyHouseholds(context.Background(), ownerID)

	assert.NoError(t, err)
	assert.Len(t, households, 1)
	assert.Equal(t, "The Smiths", households[0].Name)
	assert.Equal(t, ownerID, households[0].OwnerID)
	assert.Len(t, households[0].Members, 1)
}

func TestCreateHousehold(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice", Name: "Alice"}
	ak.addUser(owner)
	ownerID, _ := uuid.Parse(owner.UUID)

	dto, err := s.CreateHousehold(context.Background(), ownerID, "The Smiths")

	assert.NoError(t, err)
	assert.Equal(t, "The Smiths", dto.Name)
	assert.Equal(t, ownerID, dto.OwnerID)
	assert.Len(t, dto.Members, 1)
	assert.Equal(t, "Alice", dto.Members[0].DisplayName)
}

func TestCreateHousehold_NameRequired(t *testing.T) {
	s, _, _ := newTestService()
	_, err := s.CreateHousehold(context.Background(), uuid.New(), "")
	assert.EqualError(t, err, "name is required")
}

func TestAddMember_OwnerOnly(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	notOwner := authentik.User{PK: 2, UUID: uuid.New().String(), Username: "mallory"}
	ak.addUser(owner)
	ak.addUser(notOwner)
	household := seedHousehold(ak, "The Smiths", owner)

	notOwnerID, _ := uuid.Parse(notOwner.UUID)
	_, err := s.AddMember(context.Background(), household.PK, notOwnerID, "bob")
	assert.EqualError(t, err, "forbidden")
}

func TestAddMember_Idempotent(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	bob := authentik.User{PK: 2, UUID: uuid.New().String(), Username: "bob", Name: "Bob"}
	ak.addUser(owner)
	ak.addUser(bob)
	household := seedHousehold(ak, "The Smiths", owner)
	ownerID, _ := uuid.Parse(owner.UUID)

	dto, err := s.AddMember(context.Background(), household.PK, ownerID, "bob")
	assert.NoError(t, err)
	assert.Len(t, dto.Members, 2)

	// Adding the same member again doesn't duplicate it.
	dto, err = s.AddMember(context.Background(), household.PK, ownerID, "bob")
	assert.NoError(t, err)
	assert.Len(t, dto.Members, 2)
}

func TestAddMember_UserNotFound(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	ak.addUser(owner)
	household := seedHousehold(ak, "The Smiths", owner)
	ownerID, _ := uuid.Parse(owner.UUID)

	_, err := s.AddMember(context.Background(), household.PK, ownerID, "nobody")
	assert.EqualError(t, err, "user not found")
}

func TestRemoveMember(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	bob := authentik.User{PK: 2, UUID: uuid.New().String(), Username: "bob"}
	ak.addUser(owner)
	ak.addUser(bob)
	household := seedHousehold(ak, "The Smiths", owner, bob)
	ownerID, _ := uuid.Parse(owner.UUID)
	bobID, _ := uuid.Parse(bob.UUID)

	err := s.RemoveMember(context.Background(), household.PK, ownerID, bobID)
	assert.NoError(t, err)

	g, _ := ak.GetGroup(context.Background(), household.PK)
	assert.NotContains(t, g.Users, bob.PK)
}

func TestRemoveMember_OwnerCannotBeRemoved(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	ak.addUser(owner)
	household := seedHousehold(ak, "The Smiths", owner)
	ownerID, _ := uuid.Parse(owner.UUID)

	err := s.RemoveMember(context.Background(), household.PK, ownerID, ownerID)
	assert.EqualError(t, err, "owner cannot be removed")
}

func TestLeaveHousehold(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	bob := authentik.User{PK: 2, UUID: uuid.New().String(), Username: "bob"}
	ak.addUser(owner)
	ak.addUser(bob)
	household := seedHousehold(ak, "The Smiths", owner, bob)
	bobID, _ := uuid.Parse(bob.UUID)

	err := s.LeaveHousehold(context.Background(), household.PK, bobID)
	assert.NoError(t, err)

	g, _ := ak.GetGroup(context.Background(), household.PK)
	assert.NotContains(t, g.Users, bob.PK)
}

func TestLeaveHousehold_OwnerCannotLeave(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	ak.addUser(owner)
	household := seedHousehold(ak, "The Smiths", owner)
	ownerID, _ := uuid.Parse(owner.UUID)

	err := s.LeaveHousehold(context.Background(), household.PK, ownerID)
	assert.EqualError(t, err, "owner cannot leave")
}

func TestDeleteHousehold(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	ak.addUser(owner)
	household := seedHousehold(ak, "The Smiths", owner)
	ownerID, _ := uuid.Parse(owner.UUID)

	err := s.DeleteHousehold(context.Background(), household.PK, ownerID)
	assert.NoError(t, err)

	_, err = ak.GetGroup(context.Background(), household.PK)
	assert.ErrorIs(t, err, authentik.ErrNotFound)
}

func TestDeleteHousehold_Forbidden(t *testing.T) {
	s, ak, _ := newTestService()
	owner := authentik.User{PK: 1, UUID: uuid.New().String(), Username: "alice"}
	notOwner := authentik.User{PK: 2, UUID: uuid.New().String(), Username: "mallory"}
	ak.addUser(owner)
	ak.addUser(notOwner)
	household := seedHousehold(ak, "The Smiths", owner)
	notOwnerID, _ := uuid.Parse(notOwner.UUID)

	err := s.DeleteHousehold(context.Background(), household.PK, notOwnerID)
	assert.EqualError(t, err, "forbidden")
}
