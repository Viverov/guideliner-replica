package connectors

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

func InitKafka(config *config.Config) {
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
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(config.Kafka.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{}
	for _, ti := range config.Kafka.Topics {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:              ti.Name,
			NumPartitions:      ti.NumPartitions,
			ReplicationFactor:  ti.ReplicationFactor,
		})
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
