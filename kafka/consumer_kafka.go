package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	n     := flag.Int("n", 1000, "number of messages to expect before exiting")
	topic := flag.String("topic", "bench_topic", "Kafka topic")
	addr  := flag.String("broker", "localhost:9092", "Kafka bootstrap broker")
	flag.Parse()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{*addr},
		Topic:   *topic,
	})       // sem GroupID → leitura “sem-compromisso”
	defer r.Close()

	var latencies []time.Duration
	for len(latencies) < *n {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("consume: %v", err)
		}
		sent, _ := strconv.ParseInt(string(m.Value), 10, 64)
		latencies = append(latencies, time.Duration(time.Now().UnixNano()-sent))
	}

	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	avg := sum / time.Duration(len(latencies))
	fmt.Println(avg.Microseconds()) // imprime só µs
}