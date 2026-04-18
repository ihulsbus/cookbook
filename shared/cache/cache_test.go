package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/olric-data/olric"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockDMap implements the internal dmap interface for unit testing.
type mockDMap struct {
	mock.Mock
}

func (m *mockDMap) Put(ctx context.Context, key string, value interface{}, options ...olric.PutOption) error {
	return m.Called(ctx, key, value).Error(0)
}

func (m *mockDMap) Get(ctx context.Context, key string) (*olric.GetResponse, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*olric.GetResponse), args.Error(1)
}

func (m *mockDMap) Delete(ctx context.Context, keys ...string) (int, error) {
	args := m.Called(ctx, keys)
	return args.Int(0), args.Error(1)
}

// ==================================================================================================

func TestCache_Put(t *testing.T) {
	md := new(mockDMap)
	ctx := context.Background()
	md.On("Put", ctx, "key1", "value1").Return(nil)

	c := newCache(ctx, md)
	err := c.Put("key1", "value1")

	assert.NoError(t, err)
	md.AssertCalled(t, "Put", ctx, "key1", "value1")
}

func TestCache_Put_Error(t *testing.T) {
	md := new(mockDMap)
	ctx := context.Background()
	md.On("Put", ctx, "key1", "value1").Return(errors.New("put error"))

	c := newCache(ctx, md)
	err := c.Put("key1", "value1")

	assert.Error(t, err)
}

func TestCache_Delete(t *testing.T) {
	md := new(mockDMap)
	ctx := context.Background()
	md.On("Delete", ctx, []string{"key1"}).Return(1, nil)

	c := newCache(ctx, md)
	err := c.Delete("key1")

	assert.NoError(t, err)
	md.AssertCalled(t, "Delete", ctx, []string{"key1"})
}

func TestCache_Delete_Error(t *testing.T) {
	md := new(mockDMap)
	ctx := context.Background()
	md.On("Delete", ctx, []string{"key1"}).Return(0, errors.New("delete error"))

	c := newCache(ctx, md)
	err := c.Delete("key1")

	assert.Error(t, err)
}

func TestCache_Get_Error(t *testing.T) {
	md := new(mockDMap)
	ctx := context.Background()
	md.On("Get", ctx, "missing").Return(nil, errors.New("key not found"))

	c := newCache(ctx, md)
	var dst string
	err := c.Get("missing", &dst)

	assert.Error(t, err)
}

func TestCache_Get_KeyNotFound(t *testing.T) {
	md := new(mockDMap)
	ctx := context.Background()
	md.On("Get", ctx, "key1").Return(nil, olric.ErrKeyNotFound)

	c := newCache(ctx, md)
	var dst string
	err := c.Get("key1", &dst)

	require.Error(t, err)
	assert.ErrorIs(t, err, olric.ErrKeyNotFound)
}
