package main

import (
	"context"
	c "instruction-service/internal/config"
	s "instruction-service/internal/instruction-service"
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

	s.InstructionService(ctx)

	ctx.Done()

}
