package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	brokers := flag.String("brokers", "localhost:9092", "comma-separated broker list")
	n := flag.Int("n", 10000, "number of messages to expect")
	topic := flag.String("topic", "bench_topic", "Kafka topic")
	groupID := flag.String("group", "bench-consumer-group", "consumer group ID")
	outputFile := flag.String("output", "", "output file for latency measurements (default: results/kafka/test_TIMESTAMP.txt)")
	flag.Parse()

	// Create output file
	var filename string
	if *outputFile == "" {
		os.MkdirAll("results/kafka", 0755)
		filename = fmt.Sprintf("results/kafka/test_%d.txt", time.Now().Unix())
	} else {
		filename = *outputFile
		os.MkdirAll(filepath.Dir(filename), 0755)
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("create output file: %v", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  split(*brokers),
		Topic:    *topic,
		GroupID:  *groupID,
		MinBytes: 1,
		MaxBytes: 1e6,
		// Configure for at-least-once delivery
		CommitInterval: time.Second,
	})
	defer r.Close()

	ctx := context.Background()
	messagesReceived := 0

	for messagesReceived < *n {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Fatalf("read message: %v", err)
		}

		// Business Contract: Message value contains timestamp as nanoseconds since Unix epoch
		sent, err := strconv.ParseInt(string(m.Value), 10, 64)
		if err != nil {
			log.Printf("invalid timestamp format: %s", string(m.Value))
			continue
		}

		// Calculate latency in nanoseconds and convert to microseconds
		latencyNs := time.Now().UnixNano() - sent
		latencyMicros := latencyNs / 1000

		// Write latency to file (one measurement per line)
		fmt.Fprintf(writer, "%d\n", latencyMicros)
		messagesReceived++
	}

	log.Printf("Received %d messages, latencies written to %s", messagesReceived, filename)
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
