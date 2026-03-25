package repository

import (
	"cisterna-mvp/tracking-service/models"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type TrackingRepository interface {
	SaveCoordinate(ctx context.Context, payload models.GPSPayload) error
}

type kafkaRepository struct {
	writer *kafka.Writer
}

func NewKafkaRepository(brokerAdress, topic string) *kafkaRepository {
	// producer
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokerAdress),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	return &kafkaRepository{writer: w}
}

func (kr *kafkaRepository) SaveCoordinate(ctx context.Context, payload models.GPSPayload) error {
	messageBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return kr.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(payload.TruckID),
		Value: messageBytes,
	})
}
