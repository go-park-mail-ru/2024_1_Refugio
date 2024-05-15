package main

import (
	"fmt"
	"github.com/gorilla/mux"
	gMailHand "mail/internal/pkg/gmail/delivery/http"
	"net/http"
	"time"
)

// https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=190385059984-h816t4cge5p847p533s1sftvqee6smbo.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2FgAuth&response_type=code&scope=https%3A%2F%2Fmail.google.com%2F&state=state-token"
func main() {
	settingTime()

	gMail := mux.NewRouter()
	gMail.HandleFunc("/getAuthURL", gMailHand.GetAuthURL).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gAuth", gMailHand.GoogleAuth).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gIncoming", gMailHand.GetIncoming).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gSent", gMailHand.GetSent).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/getByID/{id}", gMailHand.GetById).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gSpam", gMailHand.GetSpam).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gDrafts", gMailHand.GetDrafts).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gSend", gMailHand.Send).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gDelete/{id}", gMailHand.Delete).Methods("GET", "OPTIONS")
	// update???

	gMail.HandleFunc("/gGetLabels", gMailHand.GetAllLabels).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gGetAllEmailsInLabel/{name}", gMailHand.GetAllEmailsInLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gGetAllNameLabels/{id}", gMailHand.GetAllNameLabels).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/createLabel", gMailHand.CreateLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/deleteLabel/{id}", gMailHand.DeleteLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/updateLabel/{id}", gMailHand.UpdateLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/add_email/{id}", gMailHand.AddEmailInLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/delete_email/{id}", gMailHand.DeleteEmailInLabel).Methods("GET", "OPTIONS")
	//gMail.HandleFunc("/gSend", gMailHand.Send).Methods("GET", "OPTIONS")

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 8080), gMail)
	if err != nil {
		fmt.Println("Error when starting the server:", err)
	}
}

// settingTime setting local time on server
func settingTime() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error in location detection")
	}

	time.Local = loc
}
