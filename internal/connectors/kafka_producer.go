package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(config *config.Config, topic string) *kafkaProducer {
	p := &kafkaProducer{}
	w := &kafka.Writer{
		Addr:     kafka.TCP(fmt.Sprint(config.Kafka.Host, ":", config.Kafka.Port)),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	p.writer = w
	return p
}

func (p *kafkaProducer) WriteEvents(events []Event) error {
	messages := []kafka.Message{}

	for _, event := range events {
		value, err := json.Marshal(event)
		if err != nil {
			return err
		}
		messages = append(messages, kafka.Message{
			Key:   []byte(event.Name()),
			Value: value,
		})
	}

	err := p.writer.WriteMessages(context.Background(),
		messages...,
	)

	return err
}
