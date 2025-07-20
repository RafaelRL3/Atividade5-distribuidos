package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	brokers := flag.String("brokers", "localhost:9092", "comma-separated broker list")
	n := flag.Int("n", 10000, "number of messages to publish")
	topic := flag.String("topic", "bench_topic", "Kafka topic")
	flag.Parse()

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  split(*brokers),
		Topic:    *topic,
		Balancer: &kafka.LeastBytes{},
		// Configure for at-least-once delivery
		RequiredAcks: kafka.RequireOne,
		Async:        false,
	})
	defer w.Close()

	ctx := context.Background()
	for i := 0; i < *n; i++ {
		// Business Contract: Send current timestamp as nanoseconds since Unix epoch
		ts := time.Now().UnixNano()
		msg := kafka.Message{
			Key:   nil,
			Value: []byte(fmt.Sprintf("%d", ts)),
		}

		if err := w.WriteMessages(ctx, msg); err != nil {
			log.Fatalf("write message %d: %v", i, err)
		}
	}

	log.Printf("Published %d messages to Kafka (timestamp format: nanoseconds since Unix epoch)", *n)
}

func split(s string) []string {
	var a []string
	for _, p := range strings.Split(s, ",") {
		if t := strings.TrimSpace(p); t != "" {
			a = append(a, t)
		}
	}
	return a
}
