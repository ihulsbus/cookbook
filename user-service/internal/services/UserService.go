package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"user-service/internal/authentik"
	m "user-service/internal/models"
)

// householdsParentGroupName is the top-level Authentik group every
// household is created under, keeping them easy to tell apart from
// unrelated groups (e.g. cookbook-admin) a user might also belong to.
const householdsParentGroupName = "households"

// AuthentikClient is the subset of authentik.Client's methods the
// service layer needs, declared here so tests can supply a fake without
// talking to a real Authentik instance.
type AuthentikClient interface {
	FindUserByUUID(ctx context.Context, userUUID string) (authentik.User, error)
	FindUserExact(ctx context.Context, query string) (authentik.User, error)
	GetUserByPK(ctx context.Context, pk int) (authentik.User, error)
	FindGroupByName(ctx context.Context, name, parentPK string) (authentik.Group, error)
	GetGroup(ctx context.Context, pk string) (authentik.Group, error)
	CreateGroup(ctx context.Context, g authentik.Group) (authentik.Group, error)
	UpdateGroupUsers(ctx context.Context, pk string, userPKs []int) error
	DeleteGroup(ctx context.Context, pk string) error
}

// UserRepository caches display names locally so the rest of the
// system isn't calling Authentik on every request.
type UserRepository interface {
	GetDisplayName(id uuid.UUID) (string, error)
	UpsertDisplayName(id uuid.UUID, displayName string) error
}

type UserService struct {
	authentik AuthentikClient
	repo      UserRepository
	logger    *log.Logger
}

func NewUserService(authentikClient AuthentikClient, repo UserRepository, logger *log.Logger) *UserService {
	return &UserService{
		authentik: authentikClient,
		repo:      repo,
		logger:    logger,
	}
}

// cacheDisplayName best-effort upserts id's display name into the local
// cache. Logged, not returned - a cache-write failure shouldn't fail the
// read it's piggybacking on.
func (s *UserService) cacheDisplayName(id uuid.UUID, displayName string) {
	if err := s.repo.UpsertDisplayName(id, displayName); err != nil {
		s.logger.Warnf("failed to cache display name for %s: %v", id, err)
	}
}

// toUserDTO converts an Authentik user to a UserDTO. Falls back to the
// username when Authentik's `name` field is blank (it's optional there).
func toUserDTO(u authentik.User) (m.UserDTO, error) {
	id, err := uuid.Parse(u.UUID)
	if err != nil {
		return m.UserDTO{}, fmt.Errorf("authentik user %q has a non-UUID id: %w", u.Username, err)
	}

	name := u.Name
	if name == "" {
		name = u.Username
	}

	return m.UserDTO{ID: id, Username: u.Username, DisplayName: name}, nil
}

// SearchUser resolves query to exactly one user by an exact username
// or email match - never partial, so this endpoint can't be used to
// enumerate users.
func (s *UserService) SearchUser(ctx context.Context, query string) (m.UserDTO, error) {
	if query == "" {
		return m.UserDTO{}, errors.New("query is required")
	}

	au, err := s.authentik.FindUserExact(ctx, query)
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return m.UserDTO{}, errors.New("user not found")
		}
		return m.UserDTO{}, errors.New("internal server error")
	}

	dto, err := toUserDTO(au)
	if err != nil {
		s.logger.Errorf("%v", err)
		return m.UserDTO{}, errors.New("internal server error")
	}

	s.cacheDisplayName(dto.ID, dto.DisplayName)
	return dto, nil
}

// GetDisplayName resolves id to a display name (e.g. for "shared by"
// labels) - the local cache first, falling back to Authentik on a miss.
func (s *UserService) GetDisplayName(ctx context.Context, id uuid.UUID) (m.UserDTO, error) {
	if name, err := s.repo.GetDisplayName(id); err == nil {
		return m.UserDTO{ID: id, DisplayName: name}, nil
	}

	au, err := s.authentik.FindUserByUUID(ctx, id.String())
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return m.UserDTO{}, errors.New("user not found")
		}
		return m.UserDTO{}, errors.New("internal server error")
	}

	dto, err := toUserDTO(au)
	if err != nil {
		s.logger.Errorf("%v", err)
		return m.UserDTO{}, errors.New("internal server error")
	}

	s.cacheDisplayName(dto.ID, dto.DisplayName)
	return dto, nil
}

// householdsParentPK resolves the "households" parent group, creating it
// the first time it's needed. Not cached: household creation is rare
// enough that one extra Authentik call per creation is cheaper than the
// complexity of caching it correctly.
func (s *UserService) householdsParentPK(ctx context.Context) (string, error) {
	g, err := s.authentik.FindGroupByName(ctx, householdsParentGroupName, "")
	if errors.Is(err, authentik.ErrNotFound) {
		g, err = s.authentik.CreateGroup(ctx, authentik.Group{Name: householdsParentGroupName})
	}
	if err != nil {
		return "", err
	}
	return g.PK, nil
}

// isHousehold reports whether g is a household group (marked by the
// cookbook_household attribute), as opposed to an unrelated Authentik
// group (e.g. cookbook-admin) the caller happens to belong to.
func isHousehold(g authentik.Group) bool {
	v, ok := g.Attributes["cookbook_household"]
	b, _ := v.(bool)
	return ok && b
}

func householdOwnerID(g authentik.Group) (uuid.UUID, error) {
	raw, _ := g.Attributes["owner"].(string)
	return uuid.Parse(raw)
}

// toHouseholdDTO resolves g's member PKs to display names, caching
// each. Households are expected to be small, so one Authentik call per
// member is an acceptable cost given how rarely this is read.
func (s *UserService) toHouseholdDTO(ctx context.Context, g authentik.Group) (m.HouseholdDTO, error) {
	ownerID, err := householdOwnerID(g)
	if err != nil {
		return m.HouseholdDTO{}, fmt.Errorf("household %s has an invalid owner attribute: %w", g.PK, err)
	}

	members := make([]m.HouseholdMemberDTO, 0, len(g.Users))
	for _, pk := range g.Users {
		au, err := s.authentik.GetUserByPK(ctx, pk)
		if err != nil {
			return m.HouseholdDTO{}, err
		}

		dto, err := toUserDTO(au)
		if err != nil {
			return m.HouseholdDTO{}, err
		}

		s.cacheDisplayName(dto.ID, dto.DisplayName)
		members = append(members, m.HouseholdMemberDTO{ID: dto.ID, Username: dto.Username, DisplayName: dto.DisplayName})
	}

	return m.HouseholdDTO{ID: g.PK, Name: g.Name, OwnerID: ownerID, Members: members}, nil
}

// ListMyHouseholds returns every household userID belongs to.
func (s *UserService) ListMyHouseholds(ctx context.Context, userID uuid.UUID) ([]m.HouseholdDTO, error) {
	au, err := s.authentik.FindUserByUUID(ctx, userID.String())
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("internal server error")
	}

	households := make([]m.HouseholdDTO, 0, len(au.Groups))
	for _, pk := range au.Groups {
		g, err := s.authentik.GetGroup(ctx, pk)
		if err != nil {
			return nil, errors.New("internal server error")
		}
		if !isHousehold(g) {
			continue
		}

		dto, err := s.toHouseholdDTO(ctx, g)
		if err != nil {
			s.logger.Errorf("%v", err)
			return nil, errors.New("internal server error")
		}
		households = append(households, dto)
	}

	return households, nil
}

// CreateHousehold creates a new household group owned by ownerID, who
// is also its first member.
func (s *UserService) CreateHousehold(ctx context.Context, ownerID uuid.UUID, name string) (m.HouseholdDTO, error) {
	if name == "" {
		return m.HouseholdDTO{}, errors.New("name is required")
	}

	owner, err := s.authentik.FindUserByUUID(ctx, ownerID.String())
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return m.HouseholdDTO{}, errors.New("user not found")
		}
		return m.HouseholdDTO{}, errors.New("internal server error")
	}

	parentPK, err := s.householdsParentPK(ctx)
	if err != nil {
		return m.HouseholdDTO{}, errors.New("internal server error")
	}

	created, err := s.authentik.CreateGroup(ctx, authentik.Group{
		Name:   name,
		Parent: parentPK,
		Users:  []int{owner.PK},
		Attributes: map[string]any{
			"cookbook_household": true,
			"owner":              ownerID.String(),
		},
	})
	if err != nil {
		return m.HouseholdDTO{}, errors.New("internal server error")
	}

	dto, err := s.toHouseholdDTO(ctx, created)
	if err != nil {
		s.logger.Errorf("%v", err)
		return m.HouseholdDTO{}, errors.New("internal server error")
	}
	return dto, nil
}

// getOwnedHousehold fetches the household at householdID and verifies
// it's owned by requesterID. Returns "not found" if it isn't a
// household (or doesn't exist), "forbidden" if it is but requesterID
// isn't the owner.
func (s *UserService) getOwnedHousehold(ctx context.Context, householdID string, requesterID uuid.UUID) (authentik.Group, error) {
	g, err := s.getHousehold(ctx, householdID)
	if err != nil {
		return authentik.Group{}, err
	}

	ownerID, err := householdOwnerID(g)
	if err != nil {
		s.logger.Errorf("household %s has an invalid owner attribute: %v", householdID, err)
		return authentik.Group{}, errors.New("internal server error")
	}
	if ownerID != requesterID {
		return authentik.Group{}, errors.New("forbidden")
	}

	return g, nil
}

// getHousehold is getOwnedHousehold without the ownership check, for
// operations any member may perform (leave).
func (s *UserService) getHousehold(ctx context.Context, householdID string) (authentik.Group, error) {
	g, err := s.authentik.GetGroup(ctx, householdID)
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return authentik.Group{}, errors.New("not found")
		}
		return authentik.Group{}, errors.New("internal server error")
	}
	if !isHousehold(g) {
		return authentik.Group{}, errors.New("not found")
	}
	return g, nil
}

// AddMember adds a user to a household by an exact username/email
// match. Owner only; idempotent if the user is already a member.
func (s *UserService) AddMember(ctx context.Context, householdID string, requesterID uuid.UUID, query string) (m.HouseholdDTO, error) {
	if query == "" {
		return m.HouseholdDTO{}, errors.New("query is required")
	}

	g, err := s.getOwnedHousehold(ctx, householdID, requesterID)
	if err != nil {
		return m.HouseholdDTO{}, err
	}

	target, err := s.authentik.FindUserExact(ctx, query)
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return m.HouseholdDTO{}, errors.New("user not found")
		}
		return m.HouseholdDTO{}, errors.New("internal server error")
	}

	if !containsInt(g.Users, target.PK) {
		updated := append(g.Users, target.PK)
		if err := s.authentik.UpdateGroupUsers(ctx, g.PK, updated); err != nil {
			return m.HouseholdDTO{}, errors.New("internal server error")
		}
		g.Users = updated
	}

	dto, err := s.toHouseholdDTO(ctx, g)
	if err != nil {
		s.logger.Errorf("%v", err)
		return m.HouseholdDTO{}, errors.New("internal server error")
	}
	return dto, nil
}

// RemoveMember removes targetUserID from a household. Owner only;
// idempotent if the target already isn't a member. The owner can't be
// removed this way - see DeleteHousehold.
func (s *UserService) RemoveMember(ctx context.Context, householdID string, requesterID, targetUserID uuid.UUID) error {
	g, err := s.getOwnedHousehold(ctx, householdID, requesterID)
	if err != nil {
		return err
	}

	if targetUserID == requesterID {
		return errors.New("owner cannot be removed")
	}

	target, err := s.authentik.FindUserByUUID(ctx, targetUserID.String())
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return nil
		}
		return errors.New("internal server error")
	}

	remaining := removeInt(g.Users, target.PK)
	if len(remaining) == len(g.Users) {
		return nil
	}

	if err := s.authentik.UpdateGroupUsers(ctx, g.PK, remaining); err != nil {
		return errors.New("internal server error")
	}
	return nil
}

// LeaveHousehold removes userID from a household they belong to. The
// owner can't leave this way - deleting the household is required
// instead, since there's no ownership-transfer endpoint.
func (s *UserService) LeaveHousehold(ctx context.Context, householdID string, userID uuid.UUID) error {
	g, err := s.getHousehold(ctx, householdID)
	if err != nil {
		return err
	}

	ownerID, err := householdOwnerID(g)
	if err != nil {
		s.logger.Errorf("household %s has an invalid owner attribute: %v", householdID, err)
		return errors.New("internal server error")
	}
	if ownerID == userID {
		return errors.New("owner cannot leave")
	}

	me, err := s.authentik.FindUserByUUID(ctx, userID.String())
	if err != nil {
		if errors.Is(err, authentik.ErrNotFound) {
			return nil
		}
		return errors.New("internal server error")
	}

	remaining := removeInt(g.Users, me.PK)
	if len(remaining) == len(g.Users) {
		return nil
	}

	if err := s.authentik.UpdateGroupUsers(ctx, g.PK, remaining); err != nil {
		return errors.New("internal server error")
	}
	return nil
}

// DeleteHousehold deletes a household outright. Owner only.
func (s *UserService) DeleteHousehold(ctx context.Context, householdID string, requesterID uuid.UUID) error {
	g, err := s.getOwnedHousehold(ctx, householdID, requesterID)
	if err != nil {
		return err
	}

	if err := s.authentik.DeleteGroup(ctx, g.PK); err != nil {
		return errors.New("internal server error")
	}
	return nil
}

func containsInt(haystack []int, needle int) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

func removeInt(haystack []int, needle int) []int {
	out := make([]int, 0, len(haystack))
	for _, v := range haystack {
		if v != needle {
			out = append(out, v)
		}
	}
	return out
}
