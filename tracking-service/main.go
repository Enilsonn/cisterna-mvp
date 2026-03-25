package trackingservice

import (
	"cisterna-mvp/tracking-service/handler"
	"cisterna-mvp/tracking-service/repository"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	kafkaRepository := repository.NewKafkaRepository("localhost:9092", "truck_coordinates")
	gps_handler := handler.NewGPSHandler(kafkaRepository)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("api/v1/gps", gps_handler.ReciveGPS)

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
