package http

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	apiModels "mail/internal/models/delivery_models"
	"mail/internal/models/response"
	"net/http"
	"strings"
)

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
