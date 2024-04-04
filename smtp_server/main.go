package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Println("Error reading message:", err)
		return err
	}

	fmt.Println(">-------------------------------------------------<")

	// Проверяем, что письмо отправлено с домена gmail.com и адресовано нашему домену mailhub.su
	for _, recipient := range to {
		// Печатаем содержимое письма в консоль
		fmt.Println("Received mail from:", from)
		fmt.Println("To:", recipient)
		fmt.Println("Subject:", msg.Header.Get("Subject"))

		// Чтобы вывести содержимое письма, используем msg.Body
		body, err := ioutil.ReadAll(msg.Body)
		if err != nil {
			log.Println("Error reading message body:", err)
			return err
		}
		fmt.Println("Body:", string(body))

		return nil
	}

	return nil
}

func main() {
	serverAddr := "0.0.0.0:587" // Слушаем все интерфейсы на порту 587

	// Запускаем SMTP сервер
	err := smtpd.ListenAndServe(serverAddr, mailHandler, "MailHubSMTP", "")
	if err != nil {
		log.Fatal("Error starting SMTP server:", err)
	}
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
