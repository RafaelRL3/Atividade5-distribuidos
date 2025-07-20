package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	n := flag.Int("n", 10000, "number of messages to expect before exiting")
	amqpURL := flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP connection URL")
	queueName := flag.String("queue", "bench_queue", "queue name")
	outputFile := flag.String("output", "", "output file for latency measurements (default: results/rabbitmq/test_TIMESTAMP.txt)")
	flag.Parse()

	// Create output file
	var filename string
	if *outputFile == "" {
		os.MkdirAll("results/rabbitmq", 0755)
		filename = fmt.Sprintf("results/rabbitmq/test_%d.txt", time.Now().Unix())
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

	conn, err := amqp.Dial(*amqpURL)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(*queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare: %v", err)
	}

	// Set QoS to ensure at-least-once delivery
	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("qos: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("consume: %v", err)
	}

	messagesReceived := 0
	for msg := range msgs {
		if messagesReceived >= *n {
			break
		}

		// Business Contract: Message body contains timestamp as nanoseconds since Unix epoch
		sent, err := strconv.ParseInt(string(msg.Body), 10, 64)
		if err != nil {
			log.Printf("invalid timestamp format: %s", string(msg.Body))
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
