package main

import (
	"cisterna-mvp/core-service/internal/messaging"
	"cisterna-mvp/core-service/internal/repository"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	pgConnStr := "postgres://admin:adminpassword@localhost:5432/pipeiros_db?sslmode=disable"
	kafkaBrokers := []string{"localhost:9092"}
	kafkaTopic := "truck_coordinates"
	kafkaGroupID := "core-service-group"

	repo, err := repository.NewPostgresRepo(pgConnStr)
	if err != nil {
		log.Fatalf("erro fatal ao iniciar a infraestrutura de dados: %v\n", err)
	}

	consumer := messaging.NewKafkaConsumer(
		kafkaBrokers,
		kafkaTopic,
		kafkaGroupID,
		repo,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // propaga o cancelamento

	go func() {
		consumer.StartConsumer(ctx) // iniciamos o consumidor em paralelo
	}()

	// logica para encerrar a main com o consumer com control c
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	consumer.Close()
}
