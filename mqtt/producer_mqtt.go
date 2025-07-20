package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	n := flag.Int("n", 1000, "messages")
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker URI")
	topic := flag.String("topic", "bench/topic", "publish topic")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker).SetClientID("producer")
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	defer c.Disconnect(250)

	for i := 0; i < *n; i++ {
		// Business Contract: Generate Unix timestamp in nanoseconds as string payload
		ts := time.Now().UnixNano()
		payload := fmt.Sprintf("%d", ts) // Timestamp as string (~19-20 bytes)
		token := c.Publish(*topic, 1, false, payload)
		token.Wait()
	}
	log.Printf("published %d msgs to MQTT (topic: %s)\n", *n, *topic)
}
