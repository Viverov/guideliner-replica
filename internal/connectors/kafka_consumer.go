package connectors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	defaultMinBytes = 10e3 // 10 KB
	defaultMaxBytes = 10e6 // 10 MB
	defaultMaxWait  = time.Millisecond * 100
)

type kafkaConsumer struct {
	reader        *kafka.Reader
	typeResolvers map[EventName]Event
}

type KafkaConsumerOptions struct {
	Topic         string
	ConsumerGroup string
	MinBytes      int
	MaxBytes      int
	MaxWait       time.Duration
}

func NewKafkaConsumer(config *config.Config, options KafkaConsumerOptions) *kafkaConsumer {
	c := &kafkaConsumer{}

	// Required fields
	if options.Topic == "" || options.ConsumerGroup == "" {
		panic("Topic and ConsumerGroup must be defined in KafkaConsumerOptions")
	}

	// Min-max bytes
	if options.MinBytes == 0 {
		options.MinBytes = defaultMinBytes
	}
	if options.MaxBytes == 0 {
		if options.MinBytes > defaultMaxBytes {
			options.MaxBytes = options.MinBytes * 2
		} else {
			options.MaxBytes = defaultMaxBytes
		}
	}
	if options.MaxWait == 0 {
		options.MaxWait = defaultMaxWait
	}

	c.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{fmt.Sprint(config.Kafka.Host, ":", config.Kafka.Port)},
		GroupID:  options.ConsumerGroup,
		Topic:    options.Topic,
		MinBytes: options.MinBytes, // 10KB
		MaxBytes: options.MaxBytes, // 10MB
		MaxWait:  options.MaxWait,
	})
	return c
}

func (c *kafkaConsumer) RegisterEventTypeResolvers(resolvers map[EventName]Event) {
	c.typeResolvers = resolvers
}

func (c *kafkaConsumer) GetChannel() chan Event {
	channel := make(chan Event)
	go func() {
		err := c.readEvents(channel)
		if err != nil {
			panic(err)
		}
	}()

	return channel
}

func (c *kafkaConsumer) readEvents(channel chan<- Event) error {
	for {
		m, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		e, exists := c.typeResolvers[EventName(m.Key)]
		if !exists {
			fmt.Printf("Unregistered event: %s", m.Key)
			continue
		}

		err = json.Unmarshal(m.Value, e)
		if err != nil {
			fmt.Printf("Invalid message format: %s, error: %s", m.Value, err.Error())
		}
		channel <- e
	}
}
