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

	fmt.Println("handleCommand")

	// Implement your logic for each command
	// Example:
	if strings.HasPrefix(command, "HELO") || strings.HasPrefix(command, "EHLO") {
		s.handleHELO(conn, writer, command)
	} else {
		fmt.Println("OK")
		s.sendResponse(writer, 500, "Command not recognized")
	}

	writer.Flush()
	fmt.Println("END")
}

func (s *Server) handleHELO(conn net.Conn, writer *bufio.Writer, command string) {
	s.sendResponse(writer, 250, "Hello "+s.Address)
}

func (s *Server) handleQUIT(conn net.Conn, writer *bufio.Writer) {
	// Implement QUIT command logic
}

func (s *Server) sendResponse(writer *bufio.Writer, code int, message string) {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	writer.WriteString(response)
}

func (s *Server) sendEmail(from string, to string, data string) {
	// Implement email sending logic
}
