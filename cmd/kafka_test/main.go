package main

import (
	"fmt"
	"github.com/Viverov/guideliner/internal/config"
	"github.com/Viverov/guideliner/internal/connectors"
	"os"
	"strings"
	"time"
)

func main() {
	// Setup config
	env := strings.ToUpper(os.Getenv("GUIDELINER_ENV"))
	cfg := config.InitConfig(env, "./config.json")

	prod := connectors.NewKafkaProducer(cfg, "test")

	cons := connectors.NewKafkaConsumer(cfg, connectors.KafkaConsumerOptions{
		Topic:         "test",
		ConsumerGroup: "testgroup",
		MinBytes:      10e3,
		MaxBytes:      10e6,
		MaxWait:       time.Millisecond * 50,
	})
	cons.RegisterEventTypeResolvers(map[connectors.EventName]connectors.Event{
		connectors.DumbEventName: &connectors.DumbEvent{},
	})

	go handleEvents(cons.GetChannel())
	go produceEvents(prod)

	for {
	}
}

func handleEvents(channel chan connectors.Event) {
	for event := range channel {
		fmt.Println(time.Now())
		fmt.Println(event.Name())
		switch event.(type) {
		case *connectors.DumbEvent:
			e := event.(*connectors.DumbEvent)
			fmt.Println(e.Dumb)
		default:
			fmt.Println("Undefined type")
		}
	}
}

func produceEvents(producer connectors.Producer) {
	for {
		err := producer.WriteEvents([]connectors.Event{
			&connectors.DumbEvent{
				Dumb: fmt.Sprintf("dumb value: %s", time.Now()),
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Println("Event produced")

		time.Sleep(1 * time.Second)
	}
}
