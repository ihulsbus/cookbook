package cache

import (
	"context"

	"github.com/olric-data/olric"
	"github.com/olric-data/olric/config"
	log "github.com/sirupsen/logrus"
)

func NewCacheModule(ctx context.Context, logger *log.Logger) (*olric.Olric, error) {
	c := config.New("local")

	ctx, cancel := context.WithCancel(ctx)
	c.Started = func() {
		defer cancel()
		logger.Info("cache: started and ready to accept connections")
	}

	db, err := olric.New(c)
	if err != nil {
		logger.Fatalf("cache: failed to create instance: %v", err)
	}

	go func() {
		err = db.Start()
		if err != nil {
			logger.Fatalf("cache: start returned an error: %v", err)
		}
	}()

	<-ctx.Done()

	return db, err
}
