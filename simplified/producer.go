package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    addr := flag.String("addr", "localhost:9000", "queue server <host:port>")
    n := flag.Int("n", 1000, "number of messages to publish")
    flag.Parse()

    conn, err := net.Dial("tcp", *addr)
    if err != nil {
        log.Fatalf("dial: %v", err)
    }
    defer conn.Close()
    writer := bufio.NewWriter(conn)
    reader := bufio.NewReader(conn) // to consume the "OK" acknowledgements

    for i := 0; i < *n; i++ {
        ts := time.Now().UnixNano()
        fmt.Fprintf(writer, "PUSH %d\n", ts)
        writer.Flush()
        // wait for ack so we don't overload TCP buffers (optional)
        reader.ReadString('\n')
    }
    log.Printf("published %d messages\n", *n)
}
