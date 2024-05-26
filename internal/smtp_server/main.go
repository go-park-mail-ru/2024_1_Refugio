package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/mail"
	"regexp"

	"github.com/jhillyerd/enmime"
	"github.com/mhale/smtpd"
)

const sendEmailEndpoint = "https://mailhub.su/api/v1/auth/sendOther"
const addFileEndpoint = "https://mailhub.su/api/v1/auth/addFileOther"

func main() {
	err := smtpd.ListenAndServe("0.0.0.0:587", mailHandler, "MailHubSMTP", "")
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

	var fileURLs []string
	for _, attachment := range env.Attachments {
		log.Println(attachment.FileName)
		fileURL, err := uploadFile(attachment.FileName, bytes.NewReader(attachment.Content))
		log.Println(fileURL)
		if err != nil {
			log.Printf("Ошибка при загрузке файла '%s': %v", attachment.FileName, err)
			continue
		}
		fileURLs = append(fileURLs, fileURL)
	}

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

func isValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@mailhub\.su$`)

	return emailRegex.MatchString(email)
}

func uploadFile(fileName string, file io.Reader) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", err
	}
	writer.Close()

	req, err := http.NewRequest("POST", addFileEndpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный код состояния при загрузке файла: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	fileURL, ok := response["FileId"].(string)
	if !ok {
		return "", fmt.Errorf("ошибка при получении URL файла из ответа")
	}

	return fileURL, nil
}
