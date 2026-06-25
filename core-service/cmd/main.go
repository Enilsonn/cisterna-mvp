package main

import (
	"cisterna-mvp/core-service/internal/api"
	"cisterna-mvp/core-service/internal/messaging"
	"cisterna-mvp/core-service/internal/repository"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// -- DB contigs
	pgConnStr := "postgres://admin:adminpassword@localhost:5432/pipeiros_db?sslmode=disable"
	repo, err := repository.NewPostgresRepo(pgConnStr)
	if err != nil {
		log.Fatalf("erro fatal ao iniciar a infraestrutura de dados: %v\n", err)
	}

	// repo, err := repository.NewInMemoryRepo()
	// if err != nil {
	// 		log.Fatal(err)
	//}

	// -- Kafka configs
	kafkaBrokers := []string{"localhost:9092"}
	kafkaTopic := "truck_coordinates"
	kafkaGroupID := "core-service-group"
	consumer := messaging.NewKafkaConsumer(
		kafkaBrokers,
		kafkaTopic,
		kafkaGroupID,
		repo,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // propaga o cancelamento do consumer
	go func() {
		consumer.StartConsumer(ctx) // iniciamos o consumidor em paralelo
	}()

	// -- API configs
	mapHandler := api.NewMapHandler(repo)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/cisterns", mapHandler.HandlerGetCisters)
		r.Put("/cisterns", mapHandler.HandlerCreateCistern)

		r.Get("/trucks/{truck_id}/location", mapHandler.HandlerGetTruckCurrentStatus)
	})
	s := http.Server{
		Addr:              ":8081",
		Handler:           r,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
	}
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	// logica para encerrar a main com o consumer com control c
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	consumer.Close()
}
