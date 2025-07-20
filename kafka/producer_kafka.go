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
	brokers := flag.String("brokers", "localhost:9092", "comma‑separated broker list")
	n := flag.Int("n", 1000, "number of messages")
	topic := flag.String("topic", "bench_topic", "Kafka topic")
	flag.Parse()

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  split(*brokers),
		Topic:    *topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer w.Close()

	for i := 0; i < *n; i++ {
		// Business Contract: Generate Unix timestamp in nanoseconds as string payload
		ts := time.Now().UnixNano()
		msg := kafka.Message{
			Key:   nil,
			Value: []byte(fmt.Sprintf("%d", ts)), // Timestamp as string (~19-20 bytes)
		}
		if err := w.WriteMessages(context.Background(), msg); err != nil {
			log.Fatalf("write: %v", err)
		}
	}
	log.Printf("published %d msgs to Kafka (topic: %s)\n", *n, *topic)
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
