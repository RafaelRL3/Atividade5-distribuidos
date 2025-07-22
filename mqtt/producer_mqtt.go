package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	n := flag.Int("n", 1000, "number of messages to publish")
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker URI")
	topic := flag.String("topic", "bench/topic", "MQTT topic")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("connect: %v", token.Error())
	}
	defer client.Disconnect(250)

	for i := 0; i < *n; i++ {
		ts := time.Now().UnixNano()
		payload := fmt.Sprintf("%d", ts)
		token := client.Publish(*topic, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			log.Fatalf("publish: %v", token.Error())
		}
	}
	log.Printf("published %d messages to MQTT\n", *n)
}
