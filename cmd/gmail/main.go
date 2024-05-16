package main

import (
	"fmt"
	"github.com/gorilla/mux"
	gMailHand "mail/internal/pkg/gmail/delivery/http"
	http2 "mail/internal/pkg/gmail/gmail_handler/delivery/http"
	"net/http"
	"time"
)

// https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=190385059984-h816t4cge5p847p533s1sftvqee6smbo.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2FgAuth&response_type=code&scope=https%3A%2F%2Fmail.google.com%2F&state=state-token"
func main() {
	settingTime()

	gMail := mux.NewRouter()
	gMail.HandleFunc("/getAuthURL", gMailHand.GetAuthURL).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gAuth", gMailHand.GoogleAuth).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gIncoming", http2.GetIncoming).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gSent", http2.GetSent).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/getByID/{id}", http2.GetById).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gSpam", http2.GetSpam).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gDrafts", http2.GetDrafts).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gSend", http2.Send).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gDelete/{id}", http2.Delete).Methods("GET", "OPTIONS")
	// update???

	gMail.HandleFunc("/gGetLabels", http2.GetAllLabels).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gGetAllEmailsInLabel/{name}", http2.GetAllEmailsInLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/gGetAllNameLabels/{id}", http2.GetAllNameLabels).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/createLabel", http2.CreateLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/deleteLabel/{id}", http2.DeleteLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/updateLabel/{id}", http2.UpdateLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/add_email/{id}", http2.AddEmailInLabel).Methods("GET", "OPTIONS")
	gMail.HandleFunc("/delete_email/{id}", http2.DeleteEmailInLabel).Methods("GET", "OPTIONS")
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
