package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"mail/pkg/email"
	"mail/pkg/user"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestEmailList(t *testing.T) {
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
	r := httptest.NewRequest("GET", "/emails", nil)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()

	emailHandler.List(w, r)
	var writeEmail []email.Email
	err = json.NewDecoder(w.Body).Decode(&writeEmail)
	for i := 0; i < len(writeEmail); i++ {
		if writeEmail[i] != *email.FakeEmails[i] {
			t.Error("bad values writeEmail[i] != *email.FakeEmails[i] ")
			assert.Equal(t, *email.FakeEmails[i], writeEmail[i])
		}
	}
}

func TestEmailStatusList(t *testing.T) {
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
	expectedStatusUsers := []int{401, 200}

	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/emails", nil)
		w := httptest.NewRecorder()
		if i == 1 { // http.StatusOK
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
			r.AddCookie(cookie)
		}
		emailHandler.List(w, r)
		if w.Code != expectedStatusUsers[i] {
			t.Error("status is not ok")
			assert.Equal(t, expectedStatusUsers[1], w.Code)
		}
	}
}

func TestEmailGetByID(t *testing.T) {
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
	// Add email
	for _, body := range email.FakeEmails {
		respJSON, _ := json.Marshal(body)
		r := httptest.NewRequest("POST", "/email/add", bytes.NewReader(respJSON))
		r.AddCookie(cookie)
		w := httptest.NewRecorder()
		emailHandler.Add(w, r)
	}

	r := httptest.NewRequest("GET", "/email/{id}", nil)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	for i, _ := range email.FakeEmails {
		vars := map[string]string{"id": fmt.Sprintf("%s", strconv.Itoa(i+1))}
		r = mux.SetURLVars(r, vars)
		emailHandler.GetByID(w, r)
		var mail email.Email
		err = json.NewDecoder(w.Body).Decode(&mail)
		if mail != *email.FakeEmails[i] {
			t.Error("bad values writeEmail[i] != *email.FakeEmails[i] ")
			assert.Equal(t, *email.FakeEmails[i], mail)
		}
	}
}

func TestEmailStatusGetByID(t *testing.T) {
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
	expectedStatusUsers := []int{401, 404, 200}
	var cookId string
	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/email/{id}", nil)
		w := httptest.NewRecorder()
		if i >= 1 { // http.StatusOK
			if i == 1 {
				registerUser(t, userHandler, arrBody[0])
				cook, err := loginUser(t, userHandler, arrBody[0])
				cookId = cook
				cookie := &http.Cookie{
					Name:    "session_id",
					Value:   cook,
					Expires: time.Now().Add(90 * 24 * time.Hour),
					Path:    "/",
				}
				if err != nil {
					fmt.Println(err)
					return
				}
				r.AddCookie(cookie)
			}
			if i == 2 {
				respJSON, _ := json.Marshal(email.FakeEmails[1])
				r = httptest.NewRequest("POST", "/email/add", bytes.NewReader(respJSON))
				cookie := &http.Cookie{
					Name:    "session_id",
					Value:   cookId,
					Expires: time.Now().Add(90 * 24 * time.Hour),
					Path:    "/",
				}
				r.AddCookie(cookie)
				w = httptest.NewRecorder()
				emailHandler.Add(w, r)
			}
		}
		vars := map[string]string{"id": fmt.Sprintf("%s", strconv.Itoa(1))}
		r = mux.SetURLVars(r, vars)
		emailHandler.GetByID(w, r)
		fmt.Println("wCode ", w.Code)
		if w.Code != expectedStatusUsers[i] {
			t.Error("status is not ok")
			assert.Equal(t, expectedStatusUsers[1], w.Code)
		}
	}
}

func TestEmailDelete(t *testing.T) {
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
	expectedStatusUsers := []int{401, 400, 200}
	var cookId string
	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/email/{id}", nil)
		w := httptest.NewRecorder()
		if i >= 1 { // http.StatusOK
			if i == 1 {
				registerUser(t, userHandler, arrBody[0])
				cook, err := loginUser(t, userHandler, arrBody[0])
				cookId = cook
				cookie := &http.Cookie{
					Name:    "session_id",
					Value:   cook,
					Expires: time.Now().Add(90 * 24 * time.Hour),
					Path:    "/",
				}
				if err != nil {
					fmt.Println(err)
					return
				}
				r.AddCookie(cookie)
			}
			if i == 2 {
				respJSON, _ := json.Marshal(email.FakeEmails[1])
				r = httptest.NewRequest("POST", "/email/add", bytes.NewReader(respJSON))
				cookie := &http.Cookie{
					Name:    "session_id",
					Value:   cookId,
					Expires: time.Now().Add(90 * 24 * time.Hour),
					Path:    "/",
				}
				r.AddCookie(cookie)
				w = httptest.NewRecorder()
				emailHandler.Add(w, r)
				vars := map[string]string{"id": fmt.Sprintf("%s", strconv.Itoa(1))}
				r = mux.SetURLVars(r, vars)
			}
		}
		emailHandler.Delete(w, r)
		fmt.Println("wCode ", w.Code)
		if w.Code != expectedStatusUsers[i] {
			t.Error("status is not ok")
			assert.Equal(t, expectedStatusUsers[1], w.Code)
		}
	}
}

func TestEmailUpdate(t *testing.T) {
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
	expectedStatusUsers := []int{401, 400, 200}
	var cookId string
	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/email/update/{id}", nil)
		w := httptest.NewRecorder()
		if i >= 1 { // http.StatusOK
			if i == 1 {
				registerUser(t, userHandler, arrBody[0])
				cook, err := loginUser(t, userHandler, arrBody[0])
				cookId = cook
				cookie := &http.Cookie{
					Name:    "session_id",
					Value:   cook,
					Expires: time.Now().Add(90 * 24 * time.Hour),
					Path:    "/",
				}
				if err != nil {
					fmt.Println(err)
					return
				}
				r.AddCookie(cookie)
			}
			if i == 2 {
				respJSON, _ := json.Marshal(email.FakeEmails[1])
				r = httptest.NewRequest("POST", "/email/add", bytes.NewReader(respJSON))
				cookie := &http.Cookie{
					Name:    "session_id",
					Value:   cookId,
					Expires: time.Now().Add(90 * 24 * time.Hour),
					Path:    "/",
				}
				r.AddCookie(cookie)
				w = httptest.NewRecorder()
				emailHandler.Add(w, r)
				vars := map[string]string{"id": fmt.Sprintf("%s", strconv.Itoa(1))}
				r = mux.SetURLVars(r, vars)
			}
		}
		emailHandler.Update(w, r)
		fmt.Println("wCode ", w.Code)
		if w.Code != expectedStatusUsers[i] {
			t.Error("status is not ok")
			assert.Equal(t, expectedStatusUsers[1], w.Code)
		}
	}
}
