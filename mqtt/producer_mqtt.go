package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	n := flag.Int("n", 10000, "number of messages to publish")
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker URI")
	topic := flag.String("topic", "bench/topic", "publish topic")
	flag.Parse()

	opts := mqtt.NewClientOptions().
		AddBroker(*broker).
		SetClientID("producer").
		SetCleanSession(true)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer c.Disconnect(250)

	for i := 0; i < *n; i++ {
		// Business Contract: Send current timestamp as nanoseconds since Unix epoch
		ts := time.Now().UnixNano()
		payload := fmt.Sprintf("%d", ts)

		// Publish with QoS 1 for at-least-once delivery guarantee
		token := c.Publish(*topic, 1, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Fatalf("publish message %d: %v", i, token.Error())
		}
	}

	log.Printf("Published %d messages to MQTT (timestamp format: nanoseconds since Unix epoch)", *n)
}
