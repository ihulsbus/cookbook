// Package authentik is a minimal client for Authentik's core REST API
// (api/v3/core), used by user-service to resolve users and manage
// household groups (KIT-33/KIT-34). user-service is deliberately the
// only service configured with the service-account token this client
// sends - see user-service/internal/config/config.go and KIT-74.
package authentik

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ErrNotFound is returned when a lookup finds no matching user or group.
var ErrNotFound = errors.New("authentik: not found")

// User is the subset of Authentik's core user fields this client needs.
// UUID is the stable identifier used as the OIDC `sub` claim (KIT-33's
// "Subject mode: Based on the User's UUID" decision - see
// authentik/blueprints/cookbook-oidc.yaml). PK is Authentik's own
// internal integer id, needed only for group membership mutations,
// since a group's `users` field is a list of PKs, not UUIDs.
type User struct {
	PK       int    `json:"pk"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	// Groups is the list of group PKs (UUIDs) this user belongs to, as
	// returned alongside the user by Authentik's core user endpoints.
	Groups []string `json:"groups,omitempty"`
}

// Group is the subset of Authentik's core group fields this client
// needs. PK is a UUID string (unlike User.PK).
type Group struct {
	PK         string         `json:"pk,omitempty"`
	Name       string         `json:"name"`
	Parent     string         `json:"parent,omitempty"`
	Users      []int          `json:"users,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type userList struct {
	Results []User `json:"results"`
}

type groupList struct {
	Results []Group `json:"results"`
}

// Client calls Authentik's core API with a static service-account
// bearer token.
type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

// NewClient builds a Client against the Authentik instance that issued
// oidcIssuer - its REST API always lives at the same host, under
// /api/v3/ (see cookbook/authentik/blueprints) - authenticating every
// request with token.
func NewClient(oidcIssuer, token string) (*Client, error) {
	u, err := url.Parse(oidcIssuer)
	if err != nil {
		return nil, fmt.Errorf("authentik: invalid issuer %q: %w", oidcIssuer, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("authentik: issuer %q has no scheme/host", oidcIssuer)
	}

	return &Client{
		baseURL: u.Scheme + "://" + u.Host,
		token:   token,
		http:    &http.Client{Timeout: 10 * time.Second},
	}, nil
}

// do sends one API request and decodes a JSON response into out (when
// non-nil). A 404 response is reported as ErrNotFound so callers can
// branch on it with errors.Is, same as gorm.ErrRecordNotFound elsewhere
// in this codebase.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("authentik: encoding request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	reqURL := c.baseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return fmt.Errorf("authentik: building request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("authentik: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentik: unexpected status %d from %s %s: %s", resp.StatusCode, method, path, respBody)
	}

	if out == nil {
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("authentik: decoding response: %w", err)
	}
	return nil
}

// FindUserByUUID resolves the Authentik user with the given UUID (the
// OIDC `sub` claim). Returns ErrNotFound if no such user exists.
func (c *Client) FindUserByUUID(ctx context.Context, userUUID string) (User, error) {
	var list userList
	if err := c.do(ctx, http.MethodGet, "/api/v3/core/users/", url.Values{"uuid": {userUUID}}, nil, &list); err != nil {
		return User{}, err
	}
	if len(list.Results) == 0 {
		return User{}, ErrNotFound
	}
	return list.Results[0], nil
}

// FindUserExact looks up a user by an exact username or email match only
// (KIT-33) - never a fuzzy/partial match, to avoid user enumeration via
// a search endpoint. Tries username first, then email.
func (c *Client) FindUserExact(ctx context.Context, query string) (User, error) {
	var list userList
	if err := c.do(ctx, http.MethodGet, "/api/v3/core/users/", url.Values{"username": {query}}, nil, &list); err != nil {
		return User{}, err
	}
	if len(list.Results) > 0 {
		return list.Results[0], nil
	}

	list = userList{}
	if err := c.do(ctx, http.MethodGet, "/api/v3/core/users/", url.Values{"email": {query}}, nil, &list); err != nil {
		return User{}, err
	}
	if len(list.Results) > 0 {
		return list.Results[0], nil
	}

	return User{}, ErrNotFound
}

// GetUserByPK resolves one of a group's member PKs back to a User, for
// rendering household member lists with display names.
func (c *Client) GetUserByPK(ctx context.Context, pk int) (User, error) {
	var u User
	if err := c.do(ctx, http.MethodGet, "/api/v3/core/users/"+strconv.Itoa(pk)+"/", nil, nil, &u); err != nil {
		return User{}, err
	}
	return u, nil
}

// FindGroupByName looks up a group by exact name, optionally scoped to a
// parent group's PK ("" for top-level). Returns ErrNotFound if none
// exists.
func (c *Client) FindGroupByName(ctx context.Context, name, parentPK string) (Group, error) {
	q := url.Values{"name": {name}}
	if parentPK != "" {
		q.Set("parent", parentPK)
	}

	var list groupList
	if err := c.do(ctx, http.MethodGet, "/api/v3/core/groups/", q, nil, &list); err != nil {
		return Group{}, err
	}
	if len(list.Results) == 0 {
		return Group{}, ErrNotFound
	}
	return list.Results[0], nil
}

// GetGroup fetches one group by PK.
func (c *Client) GetGroup(ctx context.Context, pk string) (Group, error) {
	var g Group
	if err := c.do(ctx, http.MethodGet, "/api/v3/core/groups/"+pk+"/", nil, nil, &g); err != nil {
		return Group{}, err
	}
	return g, nil
}

// CreateGroup creates a group (POST /api/v3/core/groups/) and returns it
// with its server-assigned PK.
func (c *Client) CreateGroup(ctx context.Context, g Group) (Group, error) {
	var created Group
	if err := c.do(ctx, http.MethodPost, "/api/v3/core/groups/", nil, g, &created); err != nil {
		return Group{}, err
	}
	return created, nil
}

// UpdateGroupUsers replaces a group's member list (KIT-34's
// AddMember/RemoveMember/LeaveHousehold - each reads the group's current
// Users, adjusts it, and writes the whole list back).
func (c *Client) UpdateGroupUsers(ctx context.Context, pk string, userPKs []int) error {
	return c.do(ctx, http.MethodPatch, "/api/v3/core/groups/"+pk+"/", nil, map[string]any{"users": userPKs}, nil)
}

// DeleteGroup deletes a group outright (KIT-34's DELETE /households/:id).
func (c *Client) DeleteGroup(ctx context.Context, pk string) error {
	return c.do(ctx, http.MethodDelete, "/api/v3/core/groups/"+pk+"/", nil, nil, nil)
}
