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
	Message    string
	EventType  string
	Email      string
	Role       string
	UserID     string
}

type Service struct {
	readers []*kafka.Reader
}

func NewService(brokers []string, topic, groupID string) *Service {
	readers := []*kafka.Reader{
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID,
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  groupID + "-auth",
			Topic:    "auth-events",
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}

	return &Service{readers: readers}
}

func (s *Service) Start(ctx context.Context) error {
	for _, reader := range s.readers {
		go s.consume(ctx, reader)
	}

	<-ctx.Done()
	return nil
}

func parseTimestamp(raw string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, fmt.Errorf("empty timestamp")
	}
	return time.Parse(time.RFC3339, raw)
}

func (s *Service) consume(ctx context.Context, reader *kafka.Reader) {
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("error reading kafka message: %v", err)
			continue
		}

		var entry LogEntry
		entry.ReceivedAt = time.Now().UTC()

		if reader.Config().Topic == "auth-events" {
			var payload AuthEvent
			if err := json.Unmarshal(msg.Value, &payload); err != nil {
				log.Printf("invalid auth payload: %v", err)
				reader.CommitMessages(ctx, msg)
				continue
			}

			parsedTime, err := parseTimestamp(payload.Timestamp)
			if err != nil {
				parsedTime = time.Now().UTC()
			}

			entry.Source = "auth-service"
			entry.Timestamp = parsedTime
			entry.EventType = payload.EventType
			entry.Email = payload.Email
			entry.Role = payload.Role
			entry.UserID = payload.UserID
			entry.Message = fmt.Sprintf("event=%s email=%s role=%s user_id=%s", payload.EventType, payload.Email, payload.Role, payload.UserID)
		} else {
			var payload GPSMessage
			if err := json.Unmarshal(msg.Value, &payload); err != nil {
				log.Printf("invalid gps payload: %v", err)
				reader.CommitMessages(ctx, msg)
				continue
			}

			parsedTime, err := parseTimestamp(payload.Timestamp)
			if err != nil {
				parsedTime = time.Now().UTC()
			}

			entry.Source = "tracking-service"
			entry.TruckID = payload.TruckID
			entry.Latitude = payload.Latitude
			entry.Longitude = payload.Longitude
			entry.Timestamp = parsedTime
			entry.Message = fmt.Sprintf("truck=%s lat=%f lon=%f", payload.TruckID, payload.Latitude, payload.Longitude)
		}

		fmt.Printf("[logs-service] source=%s event=%s message=%s at=%s\n", entry.Source, entry.EventType, entry.Message, entry.Timestamp.Format(time.RFC3339))

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("error committing message: %v", err)
		}
	}
}
