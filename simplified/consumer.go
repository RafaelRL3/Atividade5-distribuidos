package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	addr := flag.String("addr", "localhost:9000", "queue server <host:port>")
	n := flag.Int("n", 10000, "number of messages to expect before exiting")
	outputFile := flag.String("output", "", "output file for latency measurements (default: results/simplified/test_TIMESTAMP.txt)")
	flag.Parse()

	// Create output file
	var filename string
	if *outputFile == "" {
		os.MkdirAll("results/simplified", 0755)
		filename = fmt.Sprintf("results/simplified/test_%d.txt", time.Now().Unix())
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

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)
	connWriter := bufio.NewWriter(conn)

	messagesReceived := 0
	for messagesReceived < *n {
		// Request next message
		fmt.Fprintf(connWriter, "PULL\n")
		connWriter.Flush()

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "MSG "):
			// Business Contract: Message contains timestamp as nanoseconds since Unix epoch
			tsStr := strings.TrimPrefix(line, "MSG ")
			sent, err := strconv.ParseInt(tsStr, 10, 64)
			if err != nil {
				log.Printf("invalid timestamp format: %s", tsStr)
				continue
			}

			// Calculate latency in nanoseconds and convert to microseconds
			latencyNs := time.Now().UnixNano() - sent
			latencyMicros := latencyNs / 1000

			// Write latency to file (one measurement per line)
			fmt.Fprintf(writer, "%d\n", latencyMicros)
			messagesReceived++

		case line == "EMPTY":
			time.Sleep(100 * time.Microsecond) // Brief back-off
		default:
			log.Printf("unexpected response: %q", line)
		}
	}

	log.Printf("Received %d messages, latencies written to %s", messagesReceived, filename)
}
