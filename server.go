package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/vys/go-humanize"
)

type Server struct {
	proto   string
	addr    string
	handler func(conn net.Conn)
}

func (s *Server) ListenAndGo() error {
	ln, err := net.Listen(s.proto, s.addr)
	if err != nil {
		log.Println("Failed to listen for tcp connections on address ", s.addr, " Error: ", err)
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept connection ", conn, " due to error ", err)
			continue
		}
		log.Println("Client ", conn.RemoteAddr(), " connected")
		go s.handler(conn)
	}

}

func GoRuntimeStats() {
	m := &runtime.MemStats{}
	for {
		time.Sleep(2 * time.Second)
		log.Println("# goroutines: ", runtime.NumGoroutine())
		runtime.ReadMemStats(m)
		log.Println("Memory Acquired: ", humanize.Bytes(m.Sys))
		log.Println("Memory Used    : ", humanize.Bytes(m.Alloc))
		log.Println("# malloc       : ", m.Mallocs)
		log.Println("# free         : ", m.Frees)
		log.Println("GC enabled     : ", m.EnableGC)
		log.Println("# GC           : ", m.NumGC)
		log.Println("Last GC time   : ", m.LastGC)
		log.Println("Next GC        : ", humanize.Bytes(m.NextGC))
		//runtime.GC()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	request(conn)
	response(conn)
}

func request(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Line", line)
		if i == 0 {
			m := strings.Fields(line)[0]
			fmt.Println("Methods", m)
		}
		if line == "" {
			break
		}
		i++
	}
}

func response(conn net.Conn) {
	body := `
	Hello
	`

	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(body))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, body)
}

func main() {

	go GoRuntimeStats()

	s := &Server{proto: "tcp", addr: net.JoinHostPort("127.0.0.1", "8080"), handler: handleConnection}
	s.ListenAndGo()

	log.Println("Finished execution!")
}
