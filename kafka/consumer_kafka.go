package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
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

	var latencies []time.Duration
	ctx := context.Background()
	for len(latencies) < *n {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("read: %v", err)
		}
		// Business Contract: Message payload is Unix timestamp in nanoseconds as string
		sent, _ := strconv.ParseInt(string(m.Value), 10, 64)
		lat := time.Now().UnixNano() - sent
		latencies = append(latencies, time.Duration(lat))
	}

	// Calculate and output average latency in microseconds (Business Contract)
	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	avg := sum / time.Duration(len(latencies))

	fmt.Println(avg.Microseconds()) // Output average latency in microseconds
	// Uncomment for detailed stats:
	// stats(latencies)
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

func split(s string) []string {
	var a []string
	for _, p := range strings.Split(s, ",") {
		if t := strings.TrimSpace(p); t != "" {
			a = append(a, t)
		}
	}
	return a
}
