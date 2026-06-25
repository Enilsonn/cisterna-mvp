package main

import (
	"cisterna-mvp/menagement-service/client"
	"cisterna-mvp/menagement-service/internal/api"
	"cisterna-mvp/menagement-service/internal/repository"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:admin@localhost:5432/management_db?sslmode=disable"
	}

	coreURL := os.Getenv("CORE_SERVICE_URL")
	if coreURL == "" {
		coreURL = "http://localhost:8081"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	repo, err := repository.NewPostgresRepo(dbURL)
	if err != nil {
		log.Fatalln(err)
	}

	// repo, err := repository.NewInMemorySighRepo()
	// if err != nil {
	// 		log.Fatal(err)
	//}

	coreClient := client.NewCoreClient(coreURL)

	handler := api.NewApiHandler(repo, coreClient)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP) // fundamental

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/pipeiros", func(r chi.Router) {
			r.Post("/", handler.CreatePipeiro)
			r.Put("/", handler.UpdatePipeiro)
			r.Get("/{cpf}", handler.GetPipeiroByCPF)
			r.Get("/{uuid}/trucks", handler.GetTruckByPipeiroUUID)
			r.Get("/{uuid}/deliveries", handler.GetDeliveryByPipeiroUUID)
		})

		r.Route("/trucks", func(r chi.Router) {
			r.Post("/", handler.CreateTruck)
			r.Put("/", handler.UpdateTruck)
			r.Get("/{plate}", handler.GetTruckByPlate)
			r.Get("/{uuid}/deliveries", handler.GetDeliveryByTruckUUID)
		})

		r.Route("/cisterns", func(r chi.Router) {
			r.Post("/", handler.CreateCistern)
			r.Put("/", handler.UpdateCistern)
			r.Get("/", handler.GetCisterns)
			r.Get("/{uuid}", handler.GetCisternByUUID)
		})

		r.Route("/deliveries", func(r chi.Router) {
			r.Post("/", handler.CreateDelivery)
			r.Put("/", handler.UpdateDelivery)
			r.Get("/{uuid}", handler.GetDeliveryByUUID)
		})
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}

}
