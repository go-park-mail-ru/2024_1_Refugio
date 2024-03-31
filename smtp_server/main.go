package main

import (
	"fmt"
	"log"
	"net"
)

const (
	smtpPort = "2525"
)

type Server struct {
	Address string
}

func main() {
	server := Server{
		Address: "0.0.0.0",
	}

	server.Listen()
}

func (s *Server) Listen() {
	listener, err := net.Listen("tcp", s.Address+":"+smtpPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("SMTP server listening on %s\n", s.Address+":"+smtpPort)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		fmt.Println("Listen conn: ", conn)

		go s.handleConnection(conn) //go s.handleConnection(conn)
	}
}
