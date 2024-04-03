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

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Connect from %s\n", clientAddr)

	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read command from %s: %v\n", clientAddr, err)
				s.sendResponse(writer, 500, "Error reading command")
			}
			return
		}

		command = strings.TrimRight(command, "\r\n")
		s.handleCommand(conn, writer, command)
	}
}

func (s *Server) handleCommand(conn net.Conn, writer *bufio.Writer, command string) {
	clientAddr := conn.RemoteAddr().String()
	log.Printf("Received command from %s: %s\n", clientAddr, command)

	fmt.Println(">--------------------------------------------------------------------------------<")
	fmt.Println("START")
	log.Printf("COMMAND: %v\n", command)

	switch {
	case strings.HasPrefix(command, "EHLO") || strings.HasPrefix(command, "HELO"):
		s.handleHELO(conn, writer, command)
	case strings.HasPrefix(command, "HELP"):
		s.handleHELP(conn, writer)
	case strings.HasPrefix(command, "MAIL FROM"):
		s.handleMAILFROM(conn, writer, command)
	case strings.HasPrefix(command, "RCPT TO"):
		s.handleRCPTTO(conn, writer, command)
	case strings.HasPrefix(command, "DATA"):
		if !s.handleDATA(conn, writer) {
			return
		}

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "." {
				s.endMessageContent(conn, writer)
				break
			}
			fmt.Println(line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
		}
	case strings.HasPrefix(command, "QUIT"):
		if !s.handleQUIT(conn, writer, command) {
			conn.Close()
			return
		}
	default:
		fmt.Println("I DID NOT FIND THE METHOD")
		s.handleUnknownCommand(conn, writer, command)
	}

	fmt.Println("END")
	fmt.Println(">--------------------------------------------------------------------------------<")
}

func (s *Server) sendResponse(writer *bufio.Writer, code int, message string) {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	_, err := writer.WriteString(response)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		return
	}
	if err := writer.Flush(); err != nil {
		log.Printf("Failed to flush writer: %v\n", err)
	}
}

func (s *Server) handleHELO(conn net.Conn, writer *bufio.Writer, command string) {
	parts := strings.SplitN(command, " ", 2)

	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		s.sendResponse(writer, 501, "Syntax: EHLO hostname")
	} else {
		s.sendResponse(writer, 250, fmt.Sprintf("Hello %s", parts[1]))
	}
}

func (s *Server) handleHELP(conn net.Conn, writer *bufio.Writer) {
	response := "214-Commands supported:\r\n"
	response += "214-HELO EHLO MAIL RCPT DATA QUIT HELP\r\n"
	response += "214 End of HELP info\r\n"

	_, err := writer.WriteString(response)
	if err != nil {
		log.Printf("Failed to write response: %v\n", err)
		return
	}
}

func (s *Server) handleMAILFROM(conn net.Conn, writer *bufio.Writer, command string) {
	if strings.HasPrefix(command, "MAIL FROM:") {
		sender := strings.TrimPrefix(command, "MAIL FROM:")
		sender = strings.TrimSpace(sender)

		if sender == "" {
			s.sendResponse(writer, 501, "Syntax error in parameters or arguments")
			return
		}

		fmt.Println(sender)

		s.sendResponse(writer, 250, "OK")
	} else {
		s.sendResponse(writer, 501, "Syntax error in parameters or arguments")
	}
}

func (s *Server) handleRCPTTO(conn net.Conn, writer *bufio.Writer, command string) {
	if strings.HasPrefix(command, "RCPT TO:") {
		recipient := strings.TrimPrefix(command, "RCPT TO:")
		recipient = strings.TrimSpace(recipient)

		if recipient == "" {
			s.sendResponse(writer, 501, "Syntax error in parameters or arguments")
			return
		}

		fmt.Println(recipient)

		s.sendResponse(writer, 250, "OK")
	} else {
		s.sendResponse(writer, 501, "Syntax error in parameters or arguments")
	}
}

func (s *Server) handleDATA(conn net.Conn, writer *bufio.Writer) bool {
	s.startMessageContent(conn, writer)
	return true
}

func (s *Server) startMessageContent(conn net.Conn, writer *bufio.Writer) {
	s.sendResponse(writer, 354, "Start mail input; end with <CRLF>.<CRLF>")
}

func (s *Server) endMessageContent(conn net.Conn, writer *bufio.Writer) {
	s.sendResponse(writer, 250, "Message accepted for delivery")
}

func (s *Server) handleQUIT(conn net.Conn, writer *bufio.Writer, command string) bool {
	if strings.TrimSpace(command) == "QUIT" {
		s.sendResponse(writer, 221, "Closing connection. Goodbye!")
		return false
	} else {
		s.sendResponse(writer, 500, "Syntax error in parameters or arguments")
		return true
	}
}

func (s *Server) handleUnknownCommand(conn net.Conn, writer *bufio.Writer, command string) {
	if strings.HasPrefix(command, "GET ") || strings.HasPrefix(command, "OPTIONS ") {
		s.sendResponse(writer, 502, "Error: command not implemented")
	} else {
		s.sendResponse(writer, 500, "Error: bad syntax")
	}
}

func (s *Server) sendEmail(from string, to string, data string) {
	// Implement email sending logic
}
