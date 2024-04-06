package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/mail"

	"github.com/mhale/smtpd"
)

const sendEmailEndpoint = "http://89.208.223.140:8080/email/send"

func main() {
	serverAddr := "0.0.0.0:587"

	err := smtpd.ListenAndServe(serverAddr, mailHandler, "MailHubSMTP", "")
	if err != nil {
		log.Fatal("Error starting SMTP server:", err)
	}
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Println("Error reading message:", err)
		return err
	}

	senderAddr, err := mail.ParseAddress(from)
	if err != nil {
		log.Println("Error parsing sender address:", err)
		return err
	}

	if len(to) == 0 {
		log.Println("No recipient found in 'To' header")
		return fmt.Errorf("no recipient found in 'To' header")
	}

	recipientAddr, err := mail.ParseAddress(to[0])
	if err != nil {
		log.Println("Error parsing recipient address:", err)
		return err
	}

	topic := msg.Header.Get("Subject")
	body, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Println("Error reading message body:", err)
		return err
	}

	emailData := EmailSMTP{
		Topic:          topic,
		Text:           string(body),
		SenderEmail:    senderAddr.Address,
		RecipientEmail: recipientAddr.Address,
	}

	log.Printf("Email Data: %+v", emailData)

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		log.Println("Error encoding email data to JSON:", err)
		return err
	}

	resp, err := http.Post(sendEmailEndpoint, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

/* func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Println("Error reading message:", err)
		return err
	}

	fmt.Println(">-------------------------------------------------<")

	for _, recipient := range to {
		fmt.Println("Received mail from:", from)
		fmt.Println("To:", recipient)
		fmt.Println("Subject:", msg.Header.Get("Subject"))

		body, err := ioutil.ReadAll(msg.Body)
		if err != nil {
			log.Println("Error reading message body:", err)
			return err
		}
		fmt.Println("Body:", string(body))

		// http.Post("18.90.89.76:8080/email/send")

		return nil
	}

	return nil
} */
