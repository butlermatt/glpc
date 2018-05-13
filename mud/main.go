package main

import (
	"bufio"
	"fmt"
	"github.com/butlermatt/glpc/interpreter"
	"os"
	"strings"
)

type Connection struct {
	output chan string
	input  chan []string
}

func main() {
	svr := NewServer()
	conChan := svr.Start()

	con := &Connection{output: make(chan string), input: make(chan []string)}
	conChan <- con

	inner := make(chan []string)
	go input(inner)

	for {
		select {
		case out, ok := <-con.output:
			fmt.Println(out)
			if !ok {
				break
			}
		case in := <-con.input:
			con.input <- in
		}
	}
}

func input(inChan chan<- []string) {
	r := bufio.NewScanner(os.Stdin)

	for r.Scan() {
		line := r.Text()
		if line == "" {
			continue
		}

		els := strings.Split(line, " ")
		inChan <- els
	}
}

type Server struct {
	inter *interpreter.Interpreter
	conns []*Connection
	cChan chan *Connection
}

func NewServer() *Server {
	return &Server{inter: interpreter.New()}
}

func (s *Server) Start() chan<- *Connection {
	s.cChan = make(chan *Connection)
	go s.run()
	return s.cChan
}

func (s *Server) run() {
	for {
		select {
		case conn := <-s.cChan:
			s.conns = append(s.conns, conn)
			conn.output <- "Server: Connection received.\n"
		}
	}
}
