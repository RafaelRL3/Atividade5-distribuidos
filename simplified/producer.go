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
	n := flag.Int("n", 10000, "number of messages to publish")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn) // to consume the "OK" acknowledgements

	for i := 0; i < *n; i++ {
		// Business Contract: Send current timestamp as nanoseconds since Unix epoch
		ts := time.Now().UnixNano()
		fmt.Fprintf(writer, "PUSH %d\n", ts)
		writer.Flush()

		// Wait for acknowledgment to ensure at-least-once delivery
		reader.ReadString('\n')
	}

	log.Printf("Published %d messages (timestamp format: nanoseconds since Unix epoch)", *n)
}
