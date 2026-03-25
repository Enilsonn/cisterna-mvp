package messaging

import (
	"cisterna-mvp/core-service/internal/domain"
	"cisterna-mvp/core-service/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type gpsMessageDTO struct {
	TruckID   string  `json:"truck_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type kafkaConsumer struct {
	reader *kafka.Reader
	repo   repository.PositionRepository
}

func NewKafkaConsumer(brokers []string, topic, groupID string, repo repository.PositionRepository) *kafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	return &kafkaConsumer{
		reader: r,
		repo:   repo,
	}
}

func (c *kafkaConsumer) StartConsumer(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			log.Println(err)

			// se o contexto for cancelado, matamos o processo
			if ctx.Err() != nil {
				break
			}
			continue
		}

		var dto gpsMessageDTO
		if err := json.Unmarshal(msg.Value, &dto); err != nil {
			log.Println(err)
			// -- PARA ADICIONAR: quando o dado quebrar devemos bolar alguma ideia de tratá-lo para
			// inserido novamente depois, talvez colocar em um logger para um ser humano analizar
			// como o professor Marcelo sugeriu naquele caso

			// comitamos a mensagem quebrada tirando-a da fila
			c.reader.CommitMessages(ctx, msg)
			continue
		}

		parsedTime, err := time.Parse(time.RFC3339, dto.Timestamp)
		if err != nil {
			log.Println(err)
			parsedTime = time.Now()
		}

		pos := domain.TruckPosition{
			TruckID:   dto.TruckID,
			Longitude: dto.Longitude,
			Latitude:  dto.Latitude,
			Timestamp: parsedTime,
		}

		err = c.repo.SavePosition(ctx, pos)
		if err != nil {
			log.Println(err)
			time.Sleep(2 * time.Second) // para não metralhar o banco caído (por ter dado erro)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("saved position, but error to commit to kafka: %v", err)
		} else {
			fmt.Printf("[%s] position of %s saved successfully!\n", parsedTime.Format("15:04:05"), pos.TruckID)
		}

	}
}

func (c *kafkaConsumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}
