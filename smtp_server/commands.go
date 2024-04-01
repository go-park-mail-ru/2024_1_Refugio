package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	fmt.Println("CONNECT")

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read command: %v\n", err)
			}
			return
		}

		command = strings.TrimRight(command, "\r\n")
		s.handleCommand(conn, writer, command)
	}
}

func (s *Server) handleCommand(conn net.Conn, writer *bufio.Writer, command string) {
	log.Printf("Received command: %s\n", command)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("START")
	fmt.Println("COMMAND")
	fmt.Println(command)

	switch {
	case strings.HasPrefix(command, "EHLO") || strings.HasPrefix(command, "HELO"):
		s.handleHELO(conn, writer, command)
	case strings.HasPrefix(command, "MAIL FROM"):
		s.handleMAILFROM(conn, writer, command)
	case strings.HasPrefix(command, "RCPT TO"):
		fmt.Println(">--------------------------------------------------------------------------------<")
		fmt.Println("START SEND")
		s.handleRCPTTO(conn, writer, command)
		fmt.Println(">--------------------------------------------------------------------------------<")
	case strings.HasPrefix(command, "DATA"):
		s.handleDATA(conn, writer)
	case strings.HasPrefix(command, "QUIT"):
		s.handleQUIT(conn, writer)
	default:
		fmt.Println(">--------------------------------------------------------------------------------<")
		fmt.Println("I DID NOT FIND THE METHOD")
		fmt.Println(">--------------------------------------------------------------------------------<")
		s.sendResponse(writer, 500, "Command not recognized")
	}

	fmt.Println("END")
	fmt.Println("--------------------------------------------------------------------------------")

	writer.Flush()
}

func (s *Server) handleMAILFROM(conn net.Conn, writer *bufio.Writer, command string) {
	// Process the MAIL FROM command
	from := strings.TrimPrefix(command, "MAIL FROM:")
	fmt.Println(from)
	s.sendResponse(writer, 250, "OK")
}

func (s *Server) handleDATA(conn net.Conn, writer *bufio.Writer) {
	// Process the DATA command
	s.sendResponse(writer, 354, "Start mail input; end with <CRLF>.<CRLF>")
	// Implement logic to receive email data
}

func (s *Server) handleQUIT(conn net.Conn, writer *bufio.Writer) {
	// Process the QUIT command
	s.sendResponse(writer, 221, "Goodbye")
}

func (s *Server) handleHELO(conn net.Conn, writer *bufio.Writer, command string) {
	s.sendResponse(writer, 250, "Hello "+s.Address)
}

func (s *Server) handleRCPTTO(conn net.Conn, writer *bufio.Writer, command string) {
	// Разбор команды RCPT TO для получения адреса получателя
	parts := strings.Split(command, " ")
	if len(parts) < 2 {
		s.sendResponse(writer, 501, "Syntax error in parameters or arguments")
		return
	}
	recipient := parts[2]
	fmt.Println(recipient)

	// Обработка логики для команды RCPT TO
	// Например, проверка допустимости адреса получателя и сохранение его для дальнейшей обработки
	// Отправка ответа
	s.sendResponse(writer, 250, "Recipient OK: "+recipient)
}

func (s *Server) sendResponse(writer *bufio.Writer, code int, message string) {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	writer.WriteString(response)
}

func (s *Server) sendEmail(from string, to string, data string) {
	// Implement email sending logic
}
