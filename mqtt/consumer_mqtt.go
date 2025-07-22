package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	n := flag.Int("n", 1000, "number of messages to expect before exiting")
	broker := flag.String("broker", "tcp://localhost:1883", "MQTT broker URI")
	topic := flag.String("topic", "bench/topic", "MQTT topic")
	flag.Parse()

	opts := mqtt.NewClientOptions().AddBroker(*broker)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("connect: %v", token.Error())
	}
	defer client.Disconnect(250)

	latencies := make([]time.Duration, 0, *n)
	done := make(chan struct{})

	if token := client.Subscribe(*topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
		sent, _ := strconv.ParseInt(string(msg.Payload()), 10, 64)
		lat := time.Now().UnixNano() - sent
		latencies = append(latencies, time.Duration(lat))
		if len(latencies) >= *n {
			close(done)
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatalf("subscribe: %v", token.Error())
	}

	<-done

	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	avg := sum / time.Duration(len(latencies))
	fmt.Println(avg.Microseconds()) // micro-segundos
}
