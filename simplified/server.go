package main

import (
    "bufio"
    "log"
    "net"
    "strings"
    "sync"
)

// QueueServer is a single‑instance, in‑memory FIFO queue.
type QueueServer struct {
    mu    sync.Mutex
    queue []string
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
            // PUSH <message>
            msg := strings.TrimPrefix(line, "PUSH ")
            qs.mu.Lock()
            qs.queue = append(qs.queue, msg)
            qs.mu.Unlock()
            conn.Write([]byte("OK\n"))

        case line == "PULL":
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
    log.Println("[simplified‑rabbitmq] listening on :9000 …")

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