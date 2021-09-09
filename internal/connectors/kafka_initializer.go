package connectors

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/segmentio/kafka-go"
	"net"
	"strconv"
)

func InitKafka(cfg *config.Config) {
	conn, err := kafka.Dial("tcp", fmt.Sprint(cfg.Kafka.Host, ":", cfg.Kafka.Port))
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(cfg.Kafka.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{}
	for _, ti := range []config.TopicInfo{
		cfg.Kafka.UsersTopic,
		cfg.Kafka.UsersReplyTopic,
		cfg.Kafka.GuidesTopic,
		cfg.Kafka.GuidesReplyTopic,
	} {
		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             ti.Name,
			NumPartitions:     ti.NumPartitions,
			ReplicationFactor: ti.ReplicationFactor,
		})
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}

