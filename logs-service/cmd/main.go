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
	loadEnv()

	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		topic = "truck_coordinates"
	}
	groupID := os.Getenv("LOGS_GROUP_ID")
	if groupID == "" {
		groupID = "logs-service-group"
	}

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
