package main

import (
	"context"
	c "ingredient-service/internal/config"
	s "ingredient-service/internal/ingredient-service"
	"os"
	"os/signal"
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

	s.IngredientService(ctx)

	ctx.Done()

}
