package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	n := flag.Int("n", 1000, "number of messages to publish")
	topic := flag.String("topic", "bench_topic", "Kafka topic")
	addr := flag.String("broker", "localhost:9092", "Kafka bootstrap broker")
	flag.Parse()

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{*addr},
		Topic:        *topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 0 * time.Millisecond, // ou até 0
		BatchSize:    1,                    // desativa agregação por quantidade
	})
	defer w.Close()

	for i := 0; i < *n; i++ {
		ts := time.Now().UnixNano()
		if err := w.WriteMessages(context.Background(),
			kafka.Message{Value: []byte(fmt.Sprintf("%d", ts))},
		); err != nil {
			log.Fatalf("publish: %v", err)
		}
	}
	log.Printf("published %d messages to Kafka\n", *n)
}