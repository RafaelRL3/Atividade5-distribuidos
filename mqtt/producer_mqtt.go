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
		ts := time.Now().UnixNano()
		token := c.Publish(*topic, 1, false, fmt.Sprintf("%d", ts))
		token.Wait()
	}
	log.Printf("published %d msgs to MQTT\n", *n)
}
