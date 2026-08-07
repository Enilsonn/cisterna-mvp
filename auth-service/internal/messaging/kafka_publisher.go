package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

type AuthEvent struct {
	EventType string `json:"event_type"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}

type Publisher struct {
	writer *kafka.Writer
}

func NewPublisher(brokers []string, topic string) *Publisher {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &Publisher{writer: writer}
}

func (p *Publisher) Publish(ctx context.Context, event AuthEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{Value: payload})
}

func (p *Publisher) Close() error {
	return p.writer.Close()
}

func NewAuthEvent(eventType, email, role, userID string) AuthEvent {
	return AuthEvent{
		EventType: eventType,
		Email:     email,
		Role:      role,
		UserID:    userID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
