package main

import (
	"context"
	"os"
	"os/signal"
	c "recipe-service/internal/config"
	s "recipe-service/internal/recipe-service"
	"syscall"
)

var (
	log = c.Logger
)

func main() {
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

	s.RecipeService(ctx)

	ctx.Done()

}
