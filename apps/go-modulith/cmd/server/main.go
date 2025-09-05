package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zero-one-group/go-modulith/internal/app"
)

func main() {
	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fxApp := app.NewApp()

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		fxApp.Stop(context.Background())
	}()

	log.Println("Starting application...")
	if err := fxApp.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}

	<-sigChan
	log.Println("Shutting down application...")
}
