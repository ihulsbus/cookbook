package authentik

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestServer starts an httptest server that calls handler for every
// request and asserts the bearer token sent matches token.
func newTestServer(t *testing.T, token string, handler http.HandlerFunc) (*Client, *httptest.Server) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer "+token, r.Header.Get("Authorization"))
		handler(w, r)
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(srv.URL+"/application/o/cookbook-web/", token)
	assert.NoError(t, err)

	return c, srv
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func TestNewClient_DerivesBaseURLFromIssuer(t *testing.T) {
	c, err := NewClient("https://auth.gourmedy.com/application/o/cookbook-web/", "tok")

	assert.NoError(t, err)
	assert.Equal(t, "https://auth.gourmedy.com", c.baseURL)
}

func TestNewClient_InvalidIssuer(t *testing.T) {
	_, err := NewClient("not a url", "tok")
	assert.Error(t, err)

	_, err = NewClient("/no-host", "tok")
	assert.Error(t, err)
}

func TestFindUserByUUID_Found(t *testing.T) {
	want := User{PK: 1, UUID: "11111111-1111-1111-1111-111111111111", Username: "alice"}

	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/core/users/", r.URL.Path)
		assert.Equal(t, want.UUID, r.URL.Query().Get("uuid"))
		writeJSON(w, http.StatusOK, userList{Results: []User{want}})
	})

	got, err := c.FindUserByUUID(context.Background(), want.UUID)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFindUserByUUID_NotFound(t *testing.T) {
	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, userList{})
	})

	_, err := c.FindUserByUUID(context.Background(), "missing")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestFindUserExact_UsernameMatch(t *testing.T) {
	want := User{PK: 1, Username: "alice"}

	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "alice", r.URL.Query().Get("username"))
		writeJSON(w, http.StatusOK, userList{Results: []User{want}})
	})

	got, err := c.FindUserExact(context.Background(), "alice")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFindUserExact_FallsBackToEmail(t *testing.T) {
	want := User{PK: 2, Username: "bob", Email: "bob@example.com"}
	var calls []string

	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		if u := r.URL.Query().Get("username"); u != "" {
			calls = append(calls, "username="+u)
			writeJSON(w, http.StatusOK, userList{})
			return
		}
		calls = append(calls, "email="+r.URL.Query().Get("email"))
		writeJSON(w, http.StatusOK, userList{Results: []User{want}})
	})

	got, err := c.FindUserExact(context.Background(), "bob@example.com")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, []string{"username=bob@example.com", "email=bob@example.com"}, calls)
}

func TestFindUserExact_NotFound(t *testing.T) {
	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, userList{})
	})

	_, err := c.FindUserExact(context.Background(), "nobody")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetUserByPK(t *testing.T) {
	want := User{PK: 42, Username: "carol"}

	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/core/users/42/", r.URL.Path)
		writeJSON(w, http.StatusOK, want)
	})

	got, err := c.GetUserByPK(context.Background(), 42)
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFindGroupByName(t *testing.T) {
	want := Group{PK: "g-1", Name: "households"}

	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "households", r.URL.Query().Get("name"))
		assert.Equal(t, "", r.URL.Query().Get("parent"))
		writeJSON(w, http.StatusOK, groupList{Results: []Group{want}})
	})

	got, err := c.FindGroupByName(context.Background(), "households", "")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFindGroupByName_WithParent(t *testing.T) {
	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "parent-pk", r.URL.Query().Get("parent"))
		writeJSON(w, http.StatusOK, groupList{})
	})

	_, err := c.FindGroupByName(context.Background(), "the-smiths", "parent-pk")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestGetGroup(t *testing.T) {
	want := Group{PK: "g-1", Name: "the-smiths", Users: []int{1, 2}}

	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v3/core/groups/g-1/", r.URL.Path)
		writeJSON(w, http.StatusOK, want)
	})

	got, err := c.GetGroup(context.Background(), "g-1")
	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestCreateGroup(t *testing.T) {
	input := Group{Name: "the-smiths", Parent: "parent-pk", Users: []int{1}, Attributes: map[string]any{"cookbook_household": true}}
	created := input
	created.PK = "new-pk"

	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v3/core/groups/", r.URL.Path)

		var got Group
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		assert.Equal(t, input.Name, got.Name)
		assert.Equal(t, input.Parent, got.Parent)

		writeJSON(w, http.StatusCreated, created)
	})

	got, err := c.CreateGroup(context.Background(), input)
	assert.NoError(t, err)
	assert.Equal(t, created, got)
}

func TestUpdateGroupUsers(t *testing.T) {
	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/api/v3/core/groups/g-1/", r.URL.Path)

		var body map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []any{float64(1), float64(2)}, body["users"])

		w.WriteHeader(http.StatusOK)
	})

	err := c.UpdateGroupUsers(context.Background(), "g-1", []int{1, 2})
	assert.NoError(t, err)
}

func TestDeleteGroup(t *testing.T) {
	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/v3/core/groups/g-1/", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteGroup(context.Background(), "g-1")
	assert.NoError(t, err)
}

func TestDo_UnexpectedStatus(t *testing.T) {
	c, _ := newTestServer(t, "tok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("boom"))
	})

	_, err := c.GetGroup(context.Background(), "g-1")
	assert.Error(t, err)
	assert.False(t, errors.Is(err, ErrNotFound))
}
