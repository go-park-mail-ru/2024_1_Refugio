package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/striphtmltags"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	apiModels "mail/internal/models/delivery_models"
	"mail/internal/models/response"
	"net/http"
	"os"
	"time"
)

func GetIncoming(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	req, err := srv.Users.Messages.List("me").Q("label:inbox").Do()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, _ := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		email := CreateEmailStruct(msg)
		emailsApi[i] = email
	}
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

func GetSent(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	req, err := srv.Users.Messages.List("me").Q("label:sent").Do()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		email := CreateEmailStruct(msg)
		emailsApi[i] = email
	}
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

func GetSpam(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	req, err := srv.Users.Messages.List("me").Q("label:spam").Do()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		email := CreateEmailStruct(msg)
		emailsApi[i] = email
	}

	for i, _ := range emailsApi {
		emailsApi[i].SpamStatus = true
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

func GetDrafts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	req, err := srv.Users.Messages.List("me").Q("label:drafts").Do()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, _ := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		email := CreateEmailStruct(msg)
		emailsApi[i] = email
	}

	for i, _ := range emailsApi {
		emailsApi[i].DraftStatus = true
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

func GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	fmt.Println("OK")
	msg, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		fmt.Println("Error: Unable to retrieve message: %v", err)
	}

	email := CreateEmailStruct(msg)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": email})
}

func Send(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	// input := base64.RawStdEncoding.EncodeToString([]byte("From: fedasovsergey00@gmail.com\r\nTo: fedasov03@inbox.ru\r\nSubject: My first Gmail API message\r\n\r\nЭто тестовое сообщение, отправленное через Gmail API."))
	sub := string([]rune("Привет"))
	input := base64.RawStdEncoding.EncodeToString([]byte(
		"From: fedasovsergey00@gmail.com\r\n" +
			"To: fedasov03@inbox.ru\r\n" +
			fmt.Sprintf("Subject: Hello\n\n", sub) +
			"Это тестовое сообщение, отправленное через Gmail API."))

	message := &gmail.Message{
		Raw: input,
	}

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": "OK"})
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	fmt.Println("OK")
	err = srv.Users.Messages.Delete("me", messageID).Do()
	if err != nil {
		fmt.Println("Error: Unable to retrieve message: %v", err)
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"status": "OK"})
}

func CreateEmailStruct(msg *gmail.Message) *apiModels.OtherEmail {
	email := &apiModels.OtherEmail{}
	email.ID = msg.Id
	email.DateOfDispatch = time.Unix(msg.InternalDate/1000, 0)

	fmt.Println(msg.Payload.MimeType)
	if msg.Payload.MimeType == "text/plain" {
		email = ParserMessageHeadres(email, msg)
		data, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		email.Text = striphtmltags.StripTags(string(data))
	} else if msg.Payload.MimeType == "text/html" {
		data, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		email.Text = string(data)
		email = ParserMessageHeadres(email, msg)
	} else if len(msg.Payload.Parts) != 0 {
		for _, part := range msg.Payload.Parts {
			if part.MimeType == "text/html" {
				data, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					fmt.Println("Error: ", err)
				}
				email.Text = string(data)
				email = ParserMessageHeadres(email, msg)
			}
		}
	}

	return email
}

func ParserMessageHeadres(email *apiModels.OtherEmail, msg *gmail.Message) *apiModels.OtherEmail {
	for _, mes := range msg.Payload.Headers {
		if mes.Name == "To" {
			email.RecipientEmail = mes.Value
		}
		if mes.Name == "From" {
			email.SenderEmail = mes.Value
		}
		if mes.Name == "Subject" {
			email.Topic = striphtmltags.StripTags(mes.Value)
		}
	}
	return email
}

func getClient() (*oauth2.Config, *oauth2.Token) {
	b, err := os.ReadFile("cmd/configs/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tok := &oauth2.Token{}
	data, _ := ioutil.ReadFile("token.json")
	err = json.Unmarshal(data, tok)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return config, tok
}
