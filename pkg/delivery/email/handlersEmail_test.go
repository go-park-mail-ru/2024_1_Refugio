package email

/*
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"mail/pkg/delivery/models"
	emailCore "mail/pkg/domain/models"
	"mail/pkg/repository/maps/email"
	userRepo "mail/pkg/repository/maps/user"
	"sort"

	"mail/pkg/delivery"
	"mail/pkg/delivery/session"

	converters2 "mail/pkg/delivery/converters"
	converters1 "mail/pkg/repository/converters"

	emailUc "mail/pkg/usecase/email"
	userUc "mail/pkg/usecase/user"

	userHand "mail/pkg/delivery/user"

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

var arrBody = [][]byte{
	[]byte(`{
				"id": 0,
				"login": "nasty@mail.ru",
				"name": "Nasty",
				"password": "1234",
				"surname": "Low"
			}`),
	[]byte(`{
				"id": 0,
				"login": "karpovIvan@mail.ru",
				"name": "IvAn",
				"password": "QWERTY1234",
				"surname": "Karpov"
			}`),
	[]byte(`{
				"id": 0,
				"login": "nikita@mail.ru",
				"name": "Nikita",
				"password": "qwerty1234",
				"surname": "Nosov"
			}`),
}

func registerUser(t *testing.T, userHandler *userHand.UserHandler, body []byte) error {
	r := httptest.NewRequest("POST", "/delivery/v1/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()

	userHandler.Signup(w, r)
	if w.Code != http.StatusOK {
		t.Error("status is not ok")
		return fmt.Errorf("No register")
	}

	return nil
}

func loginUser(t *testing.T, userHandler *userHand.UserHandler, body []byte) (string, error) {
	r := httptest.NewRequest("POST", "/delivery/v1/login", bytes.NewReader(body))
	w := httptest.NewRecorder()

	userHandler.Login(w, r)
	cookie := w.Header().Get("Set-Cookie")[11:43]
	if w.Code != http.StatusOK {
		assert.Equal(t, http.StatusOK, w.Code)
		t.Error("status is not ok")
		return cookie, fmt.Errorf("Not login")
	}

	return cookie, nil
}

func TestEmailAdd(t *testing.T) {
	t.Parallel()
	var (
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmptyInMemoryEmailRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
		}
	)
	expectedUsers := []map[string]uint64{{"id": 1}, {"id": 2}, {"id": 3}}

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
		r := httptest.NewRequest("POST", "/delivery/v1/email/add", bytes.NewReader(body))
		r.AddCookie(cookie)
		w := httptest.NewRecorder()

		emailHandler.Add(w, r)
		var mail map[string]interface{}
		err = json.NewDecoder(w.Body).Decode(&mail)
		if uint64(mail["body"].(map[string]interface{})["email"].(map[string]interface{})["id"].(float64)) != expectedUsers[i]["id"] {
			t.Error("status is not ok")
			assert.Equal(t, expectedUsers[i]["id"], uint64(mail["body"].(map[string]interface{})["email"].(map[string]interface{})["id"].(float64)))
			return
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
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmailMemoryRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
		}
	)
	expectedUsers := []int{200, 200, 400} // 200, 401, 400 AuthMiddleware

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
		r := httptest.NewRequest("POST", "/delivery/v1/email/add", bytes.NewReader(body))
		if i != 1 {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()

		emailHandler.Add(w, r)
		fmt.Println(w.Code, "  ", expectedUsers[i])
		if w.Code != expectedUsers[i] {
			t.Error("status is not ok")
			assert.Equal(t, expectedUsers[i], w.Code)
			return
		}
	}

}

func TestEmailList(t *testing.T) {
	t.Parallel()
	var (
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmailMemoryRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
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
	r := httptest.NewRequest("GET", "/delivery/v1/emails", nil)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()

	emailHandler.List(w, r)
	var writeEmail []models.Email
	err = json.NewDecoder(w.Body).Decode(&writeEmail)
	for i := 0; i < len(writeEmail); i++ {
		if writeEmail[i] != *converters2.EmailConvertCoreInApi(*converters1.EmailConvertDbInCore(*email.FakeEmails[uint64(i)])) {
			t.Error("bad values writeEmail[i] != *email.FakeEmails[i] ")
			assert.Equal(t, *converters2.EmailConvertCoreInApi(*converters1.EmailConvertDbInCore(*email.FakeEmails[uint64(i)])), writeEmail[i])
			return
		}
	}
}

func TestEmailStatusList(t *testing.T) {
	t.Parallel()
	var (
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmailMemoryRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
		}
	)
	expectedStatusUsers := []int{200, 200} // 401, 200 AuthMiddleware

	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/delivery/v1/emails", nil)
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
			return
		}
	}
}

type EmailBody struct {
	Email models.Email `json:"email"`
}

func TestEmailGetByID(t *testing.T) {
	t.Parallel()
	var (
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmptyInMemoryEmailRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
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
	for i := 0; i < len(email.FakeEmails); i++ {
		respJSON, _ := json.Marshal(email.FakeEmails[uint64(i)])
		r := httptest.NewRequest("POST", "/delivery/v1/email/add", bytes.NewReader(respJSON))
		r.AddCookie(cookie)
		w := httptest.NewRecorder()
		emailHandler.Add(w, r)
	}

	r := httptest.NewRequest("GET", "/delivery/v1/email/{id}", nil)
	r.AddCookie(cookie)
	w := httptest.NewRecorder()
	fakeEmail, _ := emailUseCase.GetAllEmails()
	SortEmailsByID(fakeEmail)
	for i, _ := range fakeEmail {
		vars := map[string]string{"id": fmt.Sprintf("%s", strconv.Itoa(int(i+1)))}
		r = mux.SetURLVars(r, vars)
		emailHandler.GetByID(w, r)

		var emailResponse delivery.Response
		err := json.NewDecoder(w.Body).Decode(&emailResponse)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		bodyBytes, err := json.Marshal(emailResponse.Body)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}

		bodyReader := bytes.NewReader(bodyBytes)

		var mail EmailBody
		err = json.NewDecoder(bodyReader).Decode(&mail)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		if mail.Email != *converters2.EmailConvertCoreInApi(*fakeEmail[i]) {
			t.Error("bad values writeEmail[i] != *email.FakeEmails[i] ")
			assert.Equal(t, *fakeEmail[i], mail)
			return
		}
	}
}

func TestEmailStatusGetByID(t *testing.T) {
	t.Parallel()
	var (
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmptyInMemoryEmailRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
		}
	)
	expectedStatusUsers := []int{404, 404, 200} // 401, 404, 200  AuthMiddleware
	var cookId string
	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/delivery/v1/email/{id}", nil)
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
				r = httptest.NewRequest("POST", "/delivery/v1/email/add", bytes.NewReader(respJSON))
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
			return
		}
	}
}

func TestEmailDelete(t *testing.T) {
	t.Parallel()
	var (
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmailMemoryRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
		}
	)
	expectedStatusUsers := []int{400, 400, 200} // 401, 400, 200 AuthMiddleware
	var cookId string
	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/delivery/v1/email/{id}", nil)
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
				r = httptest.NewRequest("POST", "/delivery/v1/email/add", bytes.NewReader(respJSON))
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
			return
		}
	}
}

func TestEmailUpdate(t *testing.T) {
	t.Parallel()
	var (
		sessionsManager = session.NewSessionsManager()

		emailRepository = email.NewEmailMemoryRepository()
		emailUseCase    = emailUc.NewEmailUseCase(emailRepository)
		emailHandler    = &EmailHandler{
			EmailUseCase: emailUseCase,
			Sessions:     sessionsManager,
		}

		userRepository = userRepo.NewEmptyInMemoryUserRepository()
		userUseCase    = userUc.NewUserUseCase(userRepository)
		userHandler    = &userHand.UserHandler{
			UserUseCase: userUseCase,
			Sessions:    sessionsManager,
		}
	)
	expectedStatusUsers := []int{400, 400, 200} // 401, 400, 200 AuthMiddleware
	var cookId string
	for i := 0; i < len(expectedStatusUsers); i++ {
		r := httptest.NewRequest("GET", "/delivery/v1/email/update/{id}", nil)
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
				r = httptest.NewRequest("POST", "/delivery/v1/email/add", bytes.NewReader(respJSON))
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
			return
		}
	}
}

func SortEmailsByID(emails []*emailCore.Email) {
	sort.Slice(emails, func(i, j int) bool {
		return emails[i].ID < (emails[j].ID)
	})
}
*/
