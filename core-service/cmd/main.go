package main

import (
	"cisterna-mvp/core-service/internal/repository"
	"log"
)

func main() {
	connStr := "postgres://admin:adminpassword@localhost:5432/pipeiros_db?sslmode=disable"
	_, err := repository.NewPostgresRepo(connStr)
	if err != nil {
		log.Fatalf("erro fatal ao iniciar a infraestrutura de dados: %v\n", err)
	}
}
