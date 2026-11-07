package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	c "recipe-service/internal/config"
	s "recipe-service/internal/recipe-service"
	"recipe-service/internal/replay"
	"syscall"
)

var (
	log = c.Logger

	republishAccess = flag.Bool("republish-access", false, "Replay recipe.access.changed for every recipe (including soft-deleted ones), then exit. See the runbook in README.md before running this against a live environment.")
)

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	sys := make(chan os.Signal, 1)
	signal.Notify(sys, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer func() {
		signal.Stop(sys)
		cancel()
	}()

	go func() {
		select {
		case <-sys:
			log.Info("sigterm received. Exiting.")
			cancel()
		case <-ctx.Done():
		}
	}()

	if *republishAccess {
		if err := replay.RepublishAccess(ctx, c.DatabaseClient, c.Logger); err != nil {
			log.Fatalf("republish-access failed: %v", err)
		}
		return
	}

	s.RecipeService(ctx)

	ctx.Done()

}
