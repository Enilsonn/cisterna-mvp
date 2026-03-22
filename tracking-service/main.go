package trackingservice

import (
	"cisterna-mvp/tracking-service/repository"
	"cisterna-mvp/tracking-service/service"
	"time"
)

func main() {
	kafkaRepository := repository.NewKafkaRepository("localhost:9092", "truck_coordinates")
	tracker := service.NewTrackRepository(kafkaRepository)

	tracker.ProcessAndSendLocation()

	ticker := time.NewTicker(2 * time.Minute)
	for range ticker.C {
		tracker.ProcessAndSendLocation()
	}
}
