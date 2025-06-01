package main

import (
	"context"
	"log"
	"os/signal"
	"soccer-api/internal/application"
	"syscall"
	"time"
)

func main() {
	appCtx, stopSignalListener := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignalListener()

	app, err := application.New(appCtx)
	if err != nil {
		log.Fatalf("FATAL: Application initialization failed: %v\n", err)
	}

	go func() {
		if err := app.Start(appCtx); err != nil {
			log.Fatalf("FATAL: Application start failed: %v\n", err)
		}
	}()

	<-appCtx.Done()

	cleanupCtx, cancelCleanup := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCleanup()

	if err := app.Close(cleanupCtx); err != nil {
		log.Fatalf("FATAL: Application shutdown failed: %v\n", err)
	}
}
