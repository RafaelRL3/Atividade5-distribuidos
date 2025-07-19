package main

import (
    "flag"
	"fmt"
    "log"
    "time"

    "github.com/streadway/amqp"
)

func main() {
    n := flag.Int("n", 1000, "number of messages to publish")
    flag.Parse()

    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatalf("dial: %v", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("channel: %v", err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare("bench_queue", false, false, false, false, nil)
    if err != nil {
        log.Fatalf("queue: %v", err)
    }

    for i := 0; i < *n; i++ {
        ts := time.Now().UnixNano()
        body := []byte(fmt.Sprintf("%d", ts))
        err = ch.Publish("", q.Name, false, false, amqp.Publishing{Body: body})
        if err != nil {
            log.Fatalf("publish: %v", err)
        }
    }
    log.Printf("published %d messages to RabbitMQ\n", *n)
}