package main

import (
	"cisterna-mvp/logs-service/internal/consumer"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "truck_coordinates"
	groupID := "logs-service-group"

	service := consumer.NewService(brokers, topic, groupID)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := service.Start(ctx); err != nil {
			log.Printf("logs service stopped: %v", err)
		}
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)
	<-stopCh

	log.Println("stopping logs service...")
}
