package main

import (
	"cisterna-mvp/menagement-service/client"
	"cisterna-mvp/menagement-service/internal/api"
	"cisterna-mvp/menagement-service/internal/repository"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type accessClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func requireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			if !strings.HasPrefix(authorization, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authorization, "Bearer ")
			token = strings.TrimSpace(token)
			claims, err := parseClaims(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			if claims.Role != requiredRole {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), "auth_claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseClaims(token string) (accessClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return accessClaims{}, http.ErrNoLocation
	}

	payload, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return accessClaims{}, err
	}

	var claims accessClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return accessClaims{}, err
	}

	return claims, nil
}

func main() {
	loadEnv()

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
		port = "8081"
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
	r.Use(requireRole("ADMIN"))

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
