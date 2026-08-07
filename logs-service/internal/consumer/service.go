package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type GPSMessage struct {
	TruckID   string  `json:"truck_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type LogEntry struct {
	TruckID    string
	Latitude   float64
	Longitude  float64
	Timestamp  time.Time
	ReceivedAt time.Time
	Source     string
}

type Service struct {
	reader *kafka.Reader
}

func NewService(brokers []string, topic, groupID string) *Service {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &Service{reader: reader}
}

func (s *Service) Start(ctx context.Context) error {
	for {
		msg, err := s.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			log.Printf("error reading kafka message: %v", err)
			continue
		}

		var payload GPSMessage
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("invalid payload: %v", err)
			s.reader.CommitMessages(ctx, msg)
			continue
		}

		parsedTime, err := time.Parse(time.RFC3339, payload.Timestamp)
		if err != nil {
			parsedTime = time.Now()
		}

		entry := LogEntry{
			TruckID:    payload.TruckID,
			Latitude:   payload.Latitude,
			Longitude:  payload.Longitude,
			Timestamp:  parsedTime,
			ReceivedAt: time.Now().UTC(),
			Source:     "tracking-service",
		}

		fmt.Printf("[logs-service] truck=%s lat=%f lon=%f at=%s\n", entry.TruckID, entry.Latitude, entry.Longitude, entry.Timestamp.Format(time.RFC3339))

		if err := s.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("error committing message: %v", err)
		}
	}
}
