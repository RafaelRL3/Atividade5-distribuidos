package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	n := flag.Int("n", 10000, "number of messages to publish")
	amqpURL := flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP connection URL")
	queueName := flag.String("queue", "bench_queue", "queue name")
	flag.Parse()

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

	// Declare queue with consistent settings
	q, err := ch.QueueDeclare(*queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("queue declare: %v", err)
	}

	for i := 0; i < *n; i++ {
		// Business Contract: Send current timestamp as nanoseconds since Unix epoch
		ts := time.Now().UnixNano()
		body := []byte(fmt.Sprintf("%d", ts))

		// Publish with at-least-once delivery guarantee
		err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Ensure message persistence
			Body:         body,
		})
		if err != nil {
			log.Fatalf("publish: %v", err)
		}
	}

	log.Printf("Published %d messages to RabbitMQ (timestamp format: nanoseconds since Unix epoch)", *n)
}
