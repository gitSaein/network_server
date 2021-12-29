package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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

type Header struct {
	method  string
	path    string
	version string
}

const ROOT_DIR = "."

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
	header, tf := request(conn)
	if tf {
		contents, err := openfile(header.path)
		response(conn, contents, err)
	} else {
		errorResponse(conn, "405 Method Not Allowed")
	}

}

func openfile(path string) ([]string, error) {
	file, err := os.Open(ROOT_DIR + path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, nil
}

func request(conn net.Conn) (Header, bool) {
	var header Header
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Line", line)
		if i == 0 {
			header = Header{strings.Fields(line)[0], strings.Fields(line)[1], strings.Fields(line)[2]}
		}
		if line == "" {
			break
		}
		i++
	}
	if header.method == "GET" {
		return header, true
	} else {
		return header, false
	}

}

func response(conn net.Conn, contents []string, err error) {
	var strBody string
	if err != nil {
		fmt.Fprint(conn, "HTTP/1.1 404 not found\r\n")
		strBody = "404 File Not Found"
	} else {
		fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
		strBody = strings.Join(contents, "\n")
	}

	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(strBody))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, strBody)
}

func errorResponse(conn net.Conn, msg string) {
	fmt.Fprint(conn, "HTTP/1.1 405 Method Not Allowed\r\n")
	fmt.Fprintf(conn, "Content-Length: %d\r\n", len(msg))
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, msg)
}

func main() {

	go GoRuntimeStats()

	s := &Server{proto: "tcp", addr: net.JoinHostPort("127.0.0.1", "9000"), handler: handleConnection}
	s.ListenAndGo()

	log.Println("Finished execution!")
}
