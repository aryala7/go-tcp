package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}
type Server struct {
	port   string
	ln     net.Listener
	quitch chan struct{}
	msgsch chan Message
}

func NewServer(port string) *Server {
	return &Server{
		port:   port,
		quitch: make(chan struct{}),
		msgsch: make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.ln = ln

	go s.AcceptLoop()

	<-s.quitch
	close(s.msgsch)

	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err.Error())
			continue
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr().String())
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {

	defer conn.Close()
	buff := make([]byte, 2048)

	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			continue
		}
		s.msgsch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buff[:n],
		}
	}

}

func main() {
	s := NewServer(":3000")
	go func() {
		for msg := range s.msgsch {
			fmt.Println("Message from", msg.from, ":", string(msg.payload))
		}
	}()
	log.Print(s.Start())

}
