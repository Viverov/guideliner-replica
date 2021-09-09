package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(config *config.Config, topic string) *kafkaProducer {
	p := &kafkaProducer{}
	p.createTopics(config)
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

func (p *kafkaProducer) createTopics(config *config.Config) {
	conn, err := kafka.Dial("tcp", fmt.Sprint(config.Kafka.Host, ":", config.Kafka.Port))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             "test",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
