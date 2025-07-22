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
    brokers := flag.String("brokers", "localhost:9092", "broker list")
    topic := flag.String("topic", "bench_topic", "Kafka topic to subscribe")
    n := flag.Int("n", 1000, "messages to read before exit")
    group := flag.String("group", "", "consumer group (auto if empty)")
    flag.Parse()

    // Gera GroupID único se não vier via flag — evita offsets antigos.
    if *group == "" {
        *group = fmt.Sprintf("bench-%d", time.Now().UnixNano())
    }

    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:        strings.Split(*brokers, ","),
        Topic:          *topic,
        GroupID:        *group,
        StartOffset:    kafka.LastOffset, // primeira execução começa no fim
        MinBytes:       1,
        MaxBytes:       1e6,
        CommitInterval: time.Second,      // auto‑commit offsets
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

    // estatísticas simples
    var sum time.Duration
    min := time.Duration(math.MaxInt64)
    var max time.Duration
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
    fmt.Printf("received %d msgs | min=%v max=%v avg=%v\n", len(lats), min, max, avg)
}
