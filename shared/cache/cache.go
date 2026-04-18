package cache

import (
	"context"
	"fmt"

	"github.com/olric-data/olric"
	"github.com/olric-data/olric/config"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// dmap is the subset of olric.DMap that Cache actually uses.
type dmap interface {
	Put(ctx context.Context, key string, value interface{}, options ...olric.PutOption) error
	Get(ctx context.Context, key string) (*olric.GetResponse, error)
	Delete(ctx context.Context, keys ...string) (int, error)
}

// Cache is a generic key/value cache backed by an Olric DMap.
// Callers do not need to import olric directly.
type Cache struct {
	dmap dmap
	ctx  context.Context
}

// New starts an embedded Olric instance and returns a Cache for the named map.
func New(ctx context.Context, logger Logger, name string) (*Cache, error) {
	c := config.New("local")

	startCtx, cancel := context.WithCancel(ctx)
	c.Started = func() {
		defer cancel()
		logger.Info("cache: started and ready to accept connections")
	}

	db, err := olric.New(c)
	if err != nil {
		return nil, fmt.Errorf("cache: failed to create instance: %w", err)
	}

	go func() {
		if err := db.Start(); err != nil {
			logger.Fatalf("cache: start returned an error: %v", err)
		}
	}()

	<-startCtx.Done()

	client := db.NewEmbeddedClient()
	dmap, err := client.NewDMap(name)
	if err != nil {
		return nil, fmt.Errorf("cache: failed to create dmap %q: %w", name, err)
	}

	return newCache(ctx, dmap), nil
}

// newCache constructs a Cache from an existing DMap. Used in tests.
func newCache(ctx context.Context, d dmap) *Cache {
	return &Cache{dmap: d, ctx: ctx}
}

// Put stores value under key.
func (c *Cache) Put(key string, value any) error {
	return c.dmap.Put(c.ctx, key, value)
}

// Get retrieves the value stored under key and scans it into dst.
func (c *Cache) Get(key string, dst any) error {
	gr, err := c.dmap.Get(c.ctx, key)
	if err != nil {
		return err
	}
	return gr.Scan(dst)
}

// Delete removes the entry for key.
func (c *Cache) Delete(key string) error {
	_, err := c.dmap.Delete(c.ctx, key)
	return err
}
