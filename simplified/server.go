package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
)

// QueueServer implements a simple in-memory FIFO queue with TCP protocol
// Business Contract: Handles PUSH/PULL commands for timestamp message payloads
type QueueServer struct {
	mu    sync.Mutex
	queue []string // Queue stores Unix timestamps as strings
}

func (qs *QueueServer) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return // connection closed or error
		}
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "PUSH "):
			// Business Contract: PUSH <timestamp_nanoseconds>
			// Message format: Unix timestamp in nanoseconds as string (~19-20 bytes)
			msg := strings.TrimPrefix(line, "PUSH ")
			qs.mu.Lock()
			qs.queue = append(qs.queue, msg)
			qs.mu.Unlock()
			conn.Write([]byte("OK\n"))

		case line == "PULL":
			// Business Contract: PULL returns next message or EMPTY
			qs.mu.Lock()
			if len(qs.queue) == 0 {
				qs.mu.Unlock()
				conn.Write([]byte("EMPTY\n"))
				continue
			}
			msg := qs.queue[0]
			qs.queue = qs.queue[1:]
			qs.mu.Unlock()
			conn.Write([]byte("MSG " + msg + "\n"))

		default:
			conn.Write([]byte("ERR unknown command\n"))
		}
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	defer ln.Close()
	log.Println("[custom-queue-server] listening on :9000")
	log.Println("Protocol: PUSH <timestamp> | PULL -> MSG <timestamp> | EMPTY")

	qs := &QueueServer{}
	for {
		c, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go qs.handleConn(c)
	}
}
