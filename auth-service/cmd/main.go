package main

import (
	"cisterna-mvp/auth-service/internal/api"
	"cisterna-mvp/auth-service/internal/domain"
	"cisterna-mvp/auth-service/internal/repository"
	"cisterna-mvp/auth-service/internal/service"
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	loadEnv()

	repo := repository.NewInMemoryRepository()
	authService := service.NewAuthService(repo)
	handler := api.NewHandler(authService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", handler.Register)
		r.Post("/login", handler.Login)
		r.Post("/refresh", handler.Refresh)
		r.Post("/logout", handler.Logout)
		r.Post("/validate", handler.Validate)
	})

	seedUsers(repo)

	port := getAddr("AUTH_PORT", "8082")
	log.Printf("auth-service listening on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func seedUsers(repo *repository.InMemoryRepository) {
	_ = repo.CreateUser(context.Background(), domain.User{
		ID:           "admin-1",
		Email:        "admin@cisterna.com",
		PasswordHash: service.HashPassword("admin123"),
		Role:         domain.RoleAdmin,
		IsActive:     true,
	})

	_ = repo.CreateUser(context.Background(), domain.User{
		ID:           "pipeiro-1",
		Email:        "pipeiro@cisterna.com",
		PasswordHash: service.HashPassword("pipeiro123"),
		Role:         domain.RolePipeiro,
		IsActive:     true,
	})

	_ = repo.CreateUser(context.Background(), domain.User{
		ID:           "cidadao-1",
		Email:        "cidadao@cisterna.com",
		PasswordHash: service.HashPassword("cidadao123"),
		Role:         domain.RoleCidadao,
		IsActive:     true,
	})
}
