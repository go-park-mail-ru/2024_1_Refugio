package main

import (
	"bytes"
	"log"
	"net"
	"net/mail"

	"github.com/mhale/smtpd"
)

const (
	smtpPort = "587"
)

type Server struct {
	Address string
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
	return nil
}

func main() {
	serverAddr := "0.0.0.0:587"
	smtpd.ListenAndServe(serverAddr, mailHandler, "MyServerApp", "")
}

/*
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

		go s.handleConnection(conn)
	}
}
*/
