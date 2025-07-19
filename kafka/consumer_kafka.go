package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	brokers := flag.String("brokers", "localhost:9092", "comma‑separated broker list")
	n := flag.Int("n", 1000, "number of messages")
	topic := flag.String("topic", "bench_topic", "Kafka topic")
	flag.Parse()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  split(*brokers),
		Topic:    *topic,
		MinBytes: 1, MaxBytes: 1e6,
	})
	defer r.Close()

	var lats []time.Duration
	ctx := context.Background()
	for len(lats) < *n {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("read: %v", err)
		}
		sent, _ := strconv.ParseInt(string(m.Value), 10, 64)
		lats = append(lats, time.Since(time.Unix(0, sent)))
	}

	stats(lats)
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
