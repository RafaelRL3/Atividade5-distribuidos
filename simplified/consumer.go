package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net"
    "strconv"
    "strings"
    "time"
)

func main() {
    addr := flag.String("addr", "localhost:9000", "queue server <host:port>")
    n := flag.Int("n", 1000, "number of messages to expect before exiting")
    flag.Parse()

    conn, err := net.Dial("tcp", *addr)
    if err != nil {
        log.Fatalf("dial: %v", err)
    }
    defer conn.Close()
    reader := bufio.NewReader(conn)
    writer := bufio.NewWriter(conn)

    var latencies []time.Duration
    for len(latencies) < *n {
        // request next message
        fmt.Fprintf(writer, "PULL\n")
        writer.Flush()

        line, err := reader.ReadString('\n')
        if err != nil {
            log.Fatal(err)
        }
        line = strings.TrimSpace(line)

        switch {
        case strings.HasPrefix(line, "MSG "):
            tsStr := strings.TrimPrefix(line, "MSG ")
            sent, _ := strconv.ParseInt(tsStr, 10, 64)
            lat := time.Now().UnixNano() - sent
            latencies = append(latencies, time.Duration(lat))
        case line == "EMPTY":
            time.Sleep(100 * time.Microsecond) // brief back‑off
        default:
            log.Printf("unexpected: %q", line)
        }
    }

    // simple stats
    var sum time.Duration
    for _, l := range latencies {
        sum += l
    }

    avg := sum / time.Duration(len(latencies))
    
    fmt.Println(avg.Microseconds()) // imprime só micro-segundos
    //fmt.Printf("avg=%.3fms\n", float64(avg.Nanoseconds())/1e6)
}