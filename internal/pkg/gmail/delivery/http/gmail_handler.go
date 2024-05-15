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
	"mail/internal/models/microservice_ports"
	"mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"mail/internal/pkg/utils/connect_microservice"
	"net/http"
	"os"
	"strings"
	"time"
)

// GMailHandler handles user-related HTTP requests.
type GMailHandler struct {
	Sessions domainSession.SessionsManager
}

func GoogleAuth(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")
	fmt.Println("Code: ", authCode)

	b, err := os.ReadFile("cmd/configs/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	tokFile := "token.json"
	saveToken(tokFile, tok)

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.AuthService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()
}

func GetAuthURL(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile("cmd/configs/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope) // google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println(authURL)
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"AuthURL": authURL})
}

// ---------------Email-------------

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

// ---------------OtherFolder-------------

func GetAllLabels(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	req, err := srv.Users.Labels.List("me").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	var labelsApi []*apiModels.OtherLabel
	for _, l := range req.Labels {
		if strings.Contains(l.Id, "Label") {
			label := &apiModels.OtherLabel{
				ID:   l.Id,
				Name: l.Name,
			}
			labelsApi = append(labelsApi, label)
		}
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"folders": labelsApi})
}

func GetAllEmailsInLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folderName := vars["name"]

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	req, err := srv.Users.Messages.List("me").Q(fmt.Sprintf("label:%v", folderName)).Do()
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

func GetAllNameLabels(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID := vars["id"]

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	message, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	var labelsApi []*apiModels.OtherLabel
	for _, labelId := range message.LabelIds {
		l, err := srv.Users.Labels.Get("me", labelId).Do()
		if err != nil {
			fmt.Println("Error: ", err)
		}
		if strings.Contains(l.Id, "Label") {
			label := &apiModels.OtherLabel{
				ID:   l.Id,
				Name: l.Name,
			}
			labelsApi = append(labelsApi, label)
		}
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"folders": labelsApi})

}

func CreateLabel(w http.ResponseWriter, r *http.Request) {
	// data
	name := "New Folder"

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	color := &gmail.LabelColor{
		BackgroundColor: "#000000",
		TextColor:       "#000000",
	}

	label := &gmail.Label{
		Name:  name,
		Color: color,
	}

	_, err = srv.Users.Labels.Create("me", label).Do()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": "OK"})
}

func DeleteLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	labelID := vars["id"]

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	err = srv.Users.Labels.Delete("me", labelID).Do()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": "OK"})
}

func UpdateLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	labelID := vars["id"]

	// data
	newNameLabel := "Updated Label"

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	color := &gmail.LabelColor{
		BackgroundColor: "#000000",
		TextColor:       "#000000",
	}

	newLabel := &gmail.Label{
		Name:  newNameLabel,
		Color: color,
	}

	_, err = srv.Users.Labels.Update("me", labelID, newLabel).Do()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": "OK"})
}

func AddEmailInLabel(w http.ResponseWriter, r *http.Request) {
	// data
	/*
		vars := mux.Vars(r)
		messageID := vars["messageID"]
		labelID := vars["labelID"]
	*/
	messageID := "18f7c19bdb35123a"
	labelID := []string{"Label_5002769241877771600"}

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	modifyRequest := &gmail.ModifyMessageRequest{
		AddLabelIds: labelID,
	}

	_, err = srv.Users.Messages.Modify("me", messageID, modifyRequest).Do()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": "OK"})
}

func DeleteEmailInLabel(w http.ResponseWriter, r *http.Request) {
	// data
	/*
		vars := mux.Vars(r)
		messageID := vars["messageID"]
		labelID := vars["labelID"]
	*/
	messageID := "18f7c19bdb35123a"
	labelID := []string{"Label_5002769241877771600"}

	ctx := context.Background()
	config, tok := getClient()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	modifyRequest := &gmail.ModifyMessageRequest{
		RemoveLabelIds: labelID,
	}

	_, err = srv.Users.Messages.Modify("me", messageID, modifyRequest).Do()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": "OK"})
}

// --------------------------------------

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

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
