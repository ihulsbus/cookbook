package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	c "user-service/internal/config"
	s "user-service/internal/user-service"
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

	s.UserService(ctx)

	ctx.Done()

}
