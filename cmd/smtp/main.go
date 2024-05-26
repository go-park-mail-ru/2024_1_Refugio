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
	"time"

	"github.com/jhillyerd/enmime"
	"github.com/mhale/smtpd"
)

const sendEmailEndpoint = "https://mailhub.su/api/v1/auth/sendOther"
const addFileEndpoint = "https://mailhub.su/api/v1/auth/addFileOther"
const addFileToEmailEndpoint = "https://mailhub.su/api/v1/auth/addFileToEmailOther"

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
	if decodedTopic == "" {
		decodedTopic = "Без темы"
	}

	env, err := enmime.ReadEnvelope(bytes.NewReader(data))
	if err != nil {
		log.Println("Error decoding the message body:", err)
		return err
	}
	decodedBody := env.Text
	if decodedBody == "" {
		decodedBody = "Пустое письмо"
	}

	log.Println(decodedTopic)
	log.Println(decodedBody)
	log.Println(senderAddr.Address)
	log.Println(recipientAddr.Address)

	var fileURLs []uint64
	for _, attachment := range env.Attachments {
		log.Println(attachment.FileName)
		fileURL, err := uploadFile(attachment.FileName, bytes.NewReader(attachment.Content))
		log.Println(fileURL)
		if err != nil {
			log.Printf("Error uploading the file '%s': %v", attachment.FileName, err)
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read failed:", err)
		return err
	}

	var response EmailResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		log.Println("Unmarshal failed:", err)
		return err
	}

	email := response.Body.Email
	fmt.Printf("Received email with ID: %d, Topic: %s\n", email.ID, email.Topic)

	for _, fileURL := range fileURLs {
		if err = addFileToEmail(email.ID, fileURL); err != nil {
			log.Printf("Error adding file to email: %v\n", err)
		}
	}

	return nil
}

func isValidEmailFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@mailhub\.su$`)

	return emailRegex.MatchString(email)
}

func uploadFile(fileName string, file io.Reader) (uint64, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return 0, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return 0, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", addFileEndpoint, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code when uploading a file: %d", resp.StatusCode)
	}

	var response FileResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}

	return response.Body.FileId, nil
}

func addFileToEmail(emailId uint64, fileId uint64) error {
	url := fmt.Sprintf(addFileToEmailEndpoint+"/%d/file/%d", emailId, fileId)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code when uploading a file: %d", resp.StatusCode)
	}

	return nil
}

type EmailSwag struct {
	ID             uint64    `json:"id,omitempty"`
	Topic          string    `json:"topic"`
	Text           string    `json:"text"`
	ReadStatus     bool      `json:"readStatus"`
	Flag           bool      `json:"mark,omitempty"`
	Deleted        bool      `json:"deleted"`
	DateOfDispatch time.Time `json:"dateOfDispatch,omitempty"`
	ReplyToEmailID uint64    `json:"replyToEmailId,omitempty"`
	DraftStatus    bool      `json:"draftStatus"`
	SpamStatus     bool      `json:"spamStatus"`
	SenderEmail    string    `json:"senderEmail"`
	RecipientEmail string    `json:"recipientEmail"`
}

type EmailResponse struct {
	Status int `json:"status"`
	Body   struct {
		Email EmailSwag `json:"email"`
	} `json:"body"`
}

type FileResponse struct {
	Status int `json:"status"`
	Body   struct {
		FileId uint64 `json:"FileId"`
	} `json:"body"`
}
