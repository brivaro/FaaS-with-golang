package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"worker/service"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("failed to get hostname:%v", err)
	}

	cfg := service.Config{
		NatsURL:      os.Getenv("NATS_URL"),
		ConsumerName: service.DefaultConsumerName,
		WorkerName:   fmt.Sprintf("worker_%s", hostname),
	}

	worker, err := service.NewWorker(cfg)
	if err != nil {
		log.Fatalf("failed to create worker: %v", err)
	}
	defer worker.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := worker.Setup(ctx); err != nil {
		log.Fatalf("Failed to setup worker: %v", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	if err := worker.Start(ctx); err != nil {
		log.Fatalf("Worker error: %v", err)
	}
}
