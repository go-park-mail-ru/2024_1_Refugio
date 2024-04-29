package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jhillyerd/enmime"
	"github.com/mhale/smtpd"
	"log"
	"mime"
	"net"
	"net/http"
	"net/mail"
	"regexp"
)

const sendEmailEndpoint = "https://mailhub.su/api/v1/auth/sendOther"

func main() {
	serverAddr := "0.0.0.0:587"

	err := smtpd.ListenAndServe(serverAddr, mailHandler, "MailHubSMTP", "")
	if err != nil {
		log.Fatal("Error starting SMTP server:", err)
	}
}

/*func ListenAndServe(addr string, handler smtpd.Handler, authHandler smtpd.AuthHandler) error {
	srv := &smtpd.Server{
		Addr:         addr,
		Handler:      handler,
		Appname:      "MailHubSMTP",
		Hostname:     "",
		AuthHandler:  authHandler,
		AuthRequired: true,
	}
	return srv.ListenAndServe()
}

func authHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	return true, nil
	//return string(username) == "valid" && string(password) == "password", nil
}*/

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	if true {
		return fmt.Errorf("domain in the login is not suitable")
	}

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

	if !isValidEmailFormat(recipientAddr.Address) {
		log.Println("domain in the login is not suitable:", err)
		return fmt.Errorf("domain in the login is not suitable")
	}

	topic := msg.Header.Get("Subject")
	wordDecoder := new(mime.WordDecoder)
	decodedTopic, err := wordDecoder.DecodeHeader(topic)
	if err != nil {
		log.Println("Error decoding the message subject:", err)
		return err
	}

	env, err := enmime.ReadEnvelope(bytes.NewReader(data))
	if err != nil {
		log.Println("Error decoding the message body:", err)
		return err
	}
	decodedBody := env.Text

	log.Println(decodedTopic)
	log.Println(decodedBody)
	log.Println(senderAddr.Address)
	log.Println(recipientAddr.Address)

	emailData := EmailSMTP{
		Topic:          decodedTopic,
		Text:           decodedBody,
		SenderEmail:    senderAddr.Address,
		RecipientEmail: recipientAddr.Address,
	}

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

// isValidEmailFormat checks if the provided email string matches the specific format for emails ending with "@mailhub.ru".
func isValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@mailhub\.su$`)

	return emailRegex.MatchString(email)
}
