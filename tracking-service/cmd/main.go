package main

import (
	"cisterna-mvp/tracking-service/handler"
	"cisterna-mvp/tracking-service/repository"
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
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
	kafkaRepository := repository.NewKafkaRepository("localhost:9092", "truck_coordinates")
	gps_handler := handler.NewGPSHandler(kafkaRepository)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(requireRole("PIPEIRO"))

	r.Post("/api/v1/gps", gps_handler.ReciveGPS)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,  // tempo maximo para o aparelho enviar o dado
		WriteTimeout: 30 * time.Second,  // tempo máximo para o serviço responder responder 202
		IdleTimeout:  120 * time.Second, // tempo máximo em que a conexão poderá ficar aberta
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("erro fatal no servidor")
	}
}
