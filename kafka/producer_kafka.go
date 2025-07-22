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
    topic := flag.String("topic", "bench_topic", "Kafka topic to publish to")
    n := flag.Int("n", 1000, "number of messages to publish")
    linger := flag.Duration("linger", 0, "linger/batch timeout (e.g. 1ms)")
    flag.Parse()

    w := kafka.NewWriter(kafka.WriterConfig{
        Brokers:      strings.Split(*brokers, ","),
        Topic:        *topic,
        BatchTimeout: *linger,
    })
    defer w.Close()

    ctx := context.Background()
    for i := 0; i < *n; i++ {
        ts := time.Now().UnixNano()
        if err := w.WriteMessages(ctx, kafka.Message{Value: []byte(fmt.Sprintf("%d", ts))}); err != nil {
            log.Fatalf("write: %v", err)
        }
    }
    log.Printf("published %d messages to Kafka", *n)
}
