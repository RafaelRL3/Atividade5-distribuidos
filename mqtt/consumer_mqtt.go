package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	n := flag.Int("n", 1000, "messages")
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker URI")
	topic := flag.String("topic", "bench/topic", "subscribe topic")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker).SetClientID("consumer")
	latencies := make([]time.Duration, 0, *n)

	var messageHandler mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
		// Business Contract: Message payload is Unix timestamp in nanoseconds as string
		sent, _ := strconv.ParseInt(string(m.Payload()), 10, 64)
		lat := time.Now().UnixNano() - sent
		latencies = append(latencies, time.Duration(lat))

		if len(latencies) >= *n {
			// Calculate and output average latency in microseconds (Business Contract)
			var sum time.Duration
			for _, l := range latencies {
				sum += l
			}
			avg := sum / time.Duration(len(latencies))

			fmt.Println(avg.Microseconds()) // Output average latency in microseconds
			// Uncomment for detailed stats:
			// stats(latencies)

			c.Disconnect(250)
		}
	}
	opts.SetDefaultPublishHandler(messageHandler)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	if token := c.Subscribe(*topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Block until Disconnect in handler
	select {}
}

func stats(lats []time.Duration) {
	var sum time.Duration
	min := time.Duration(math.MaxInt64)
	max := time.Duration(0)
	for _, l := range lats {
		if l < min {
			min = l
		}
		if l > max {
			max = l
		}
		sum += l
	}
	avg := sum / time.Duration(len(lats))
	fmt.Printf("received %d msgs\nmin=%v max=%v avg=%v\n", len(lats), min, max, avg)
}
