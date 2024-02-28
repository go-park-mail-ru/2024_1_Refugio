package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"mail/pkg/email"
	"mail/pkg/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var arrBodyEmail = [][]byte{
	[]byte(`{
          "dateOfDispatch": "2006-01-02T15:04:05Z",
		  "deleted": true,
		  "draftStatus": true,
		  "id": 0,
		  "mark": "string",
		  "photoId": "string",
		  "readStatus": true,
		  "replyToEmailId": 0,
		  "text": "string",
		  "topic": "string"
		}`),
	[]byte(`{
         "dateOfDispatch": "2006-01-02T15:04:05Z",
		  "deleted": true,
		  "draftStatus": true,
		  "id": 0,
		  "mark": "LALALA",
		  "photoId": "Id",
		  "readStatus": true,
		  "replyToEmailId": 0,
		  "text": "Hello",
		  "topic": "string"
		}`),
	[]byte(`{
         "dateOfDispatch": "2006-01-02T15:04:05Z",
		  "deleted": true,
		  "draftStatus": true,
		  "id": 0,
		  "mark": "lol",
		  "photoId": "lalalala",
		  "readStatus": true,
		  "replyToEmailId": 0,
		  "text": "Hi!!!",
		  "topic": "string"
		}`),
}

func TestEmailAdd(t *testing.T) {
	t.Parallel()
	var (
		emailRepository = email.NewEmailMemoryRepository()
		emailHandler    = &EmailHandler{
			EmailRepository: emailRepository,
		}

		userRepository = user.NewEmptyInMemoryUserRepository()
		userHandler    = &UserHandler{
			UserRepository: userRepository,
		}
	)
	expectedUsers := []map[string]int{{"id": 1}, {"id": 2}, {"id": 3}}

	registerUser(t, userHandler, arrBody[0])
	cook, err := loginUser(t, userHandler, arrBody[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   cook,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}

	for i, body := range arrBodyEmail {
		r := httptest.NewRequest("POST", "/email/add", bytes.NewReader(body))
		r.AddCookie(cookie)
		w := httptest.NewRecorder()

		emailHandler.Add(w, r)
		var mail map[string]int
		err = json.NewDecoder(w.Body).Decode(&mail)
		fmt.Println(mail["id"], " ", expectedUsers[i]["id"])
		if mail["id"] != expectedUsers[i]["id"] {
			t.Error("status is not ok")
			assert.Equal(t, expectedUsers[i]["id"], mail["id"])
		}
	}

}

func TestEmailStatusAdd(t *testing.T) {
	t.Parallel()
	var arrBodyEmailBadStatus = [][]byte{
		[]byte(`{
          "dateOfDispatch": "2006-01-02T15:04:05Z",
		  "deleted": true,
		  "draftStatus": true,
		  "id": 0,
		  "mark": "string",
		  "photoId": "string",
		  "readStatus": true,
		  "replyToEmailId": 0,
		  "text": "string",
		  "topic": "string"
		}`),
		[]byte(`{
         "dateOfDispatch": "2006-01-02T15:04:05Z",
		  "deleted": true,
		  "draftStatus": true,
		  "id": 0,
		  "mark": "LALALA",
		  "photoId": "Id",
		  "readStatus": true,
		  "replyToEmailId": 0,
		  "text": "Hello",
		  "topic": "string"
		}`),
		[]byte(`{
         "dateOfDispatch": "2006-01-02T15:04:05Z
		  "deleted": true,
		  "draftStatus": true,
		  "id": 0,
		  "mark": "lol",
		  "photoId": "lalalala",
		  "readStatus": true,
		  "replyToEmailId": 0,
		  "text": "Hi!!!",
		  "topic": "string"
		}`),
	}

	var (
		emailRepository = email.NewEmailMemoryRepository()
		emailHandler    = &EmailHandler{
			EmailRepository: emailRepository,
		}

		userRepository = user.NewEmptyInMemoryUserRepository()
		userHandler    = &UserHandler{
			UserRepository: userRepository,
		}
	)
	expectedUsers := []int{200, 401, 400}

	registerUser(t, userHandler, arrBody[0])
	cook, err := loginUser(t, userHandler, arrBody[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   cook,
		Expires: time.Now().Add(90 * 24 * time.Hour),
		Path:    "/",
	}

	for i, body := range arrBodyEmailBadStatus {
		r := httptest.NewRequest("POST", "/email/add", bytes.NewReader(body))
		if i != 1 {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()

		emailHandler.Add(w, r)
		fmt.Println(w.Code, "  ", expectedUsers[i])
		if w.Code != expectedUsers[i] {
			t.Error("status is not ok")
			assert.Equal(t, expectedUsers[i], w.Code)
		}
	}

}
