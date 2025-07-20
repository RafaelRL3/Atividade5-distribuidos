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
		// Business Contract: Generate Unix timestamp in nanoseconds as string payload
		ts := time.Now().UnixNano()
		fmt.Fprintf(writer, "PUSH %d\n", ts) // Send timestamp as string (~19-20 bytes)
		writer.Flush()
		// Wait for ack to avoid TCP buffer overflow
		reader.ReadString('\n')
	}
	log.Printf("published %d messages to custom queue server\n", *n)
}
