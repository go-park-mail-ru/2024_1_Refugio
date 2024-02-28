package handlers

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"mail/pkg/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

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

func registerUser(t *testing.T, userHandler *UserHandler, body []byte) error {
	r := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()

	userHandler.Signup(w, r)
	if w.Code != http.StatusOK {
		t.Error("status is not ok")
		return fmt.Errorf("No register")
	}
	return nil
}

func loginUser(t *testing.T, userHandler *UserHandler, body []byte) (string, error) {
	r := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
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

func TestSignupUser(t *testing.T) {
	t.Parallel()
	var (
		userRepository = user.NewEmptyInMemoryUserRepository()
		userHandler    = &UserHandler{
			UserRepository: userRepository,
		}
	)

	expectedUsers := []user.User{
		{
			ID:       1,
			Name:     "Nasty",
			Surname:  "Low",
			Login:    "nasty@mail.ru",
			Password: "1234",
		},
		{
			ID:       2,
			Name:     "IvAn",
			Surname:  "Karpov",
			Login:    "karpovIvan@mail.ru",
			Password: "QWERTY1234",
		},
		{
			ID:       3,
			Name:     "Nikita",
			Surname:  "Nosov",
			Login:    "nikita@mail.ru",
			Password: "qwerty1234",
		},
	}

	for _, body := range arrBody {
		r := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
		w := httptest.NewRecorder()

		userHandler.Signup(w, r)
		fmt.Println(w.Code)
		if w.Code != http.StatusOK {
			t.Error("status is not ok")
		}
	}

	allUsers, err := userHandler.UserRepository.GetAll()

	if err != nil {
		return
	}

	for i, _ := range allUsers {
		if !user.ComparingUserObjects((*allUsers[i]), expectedUsers[i]) {
			assert.Equal(t, expectedUsers[i], (*allUsers[i]))
		}
	}
}

func TestStatusSignupUser(t *testing.T) {
	t.Parallel()
	var (
		userRepository = user.NewEmptyInMemoryUserRepository()
		userHandler    = &UserHandler{
			UserRepository: userRepository,
		}
	)

	expectedStatus := []int{200, 400, 400}

	for i, body := range arrBody {
		r := httptest.NewRequest("POST", "/signup", bytes.NewReader(body))
		w := httptest.NewRecorder()

		userHandler.Signup(w, r)
		fmt.Println(w.Code)
		if w.Code != expectedStatus[i] {
			assert.Equal(t, expectedStatus[i], w.Code)
			t.Error("status is not ok")
		}
	}
}

func TestLoginUser(t *testing.T) {
	t.Parallel()
	var (
		userRepository = user.NewEmptyInMemoryUserRepository()
		userHandler    = &UserHandler{
			UserRepository: userRepository,
		}
	)
	var arrBadStatusBody = [][]byte{
		[]byte(`{
				"id": 0,
				"login": "nasty@mail.ru",
				"name": "Nasty",
				"password": "1234",
				"surname": "Low"
			}`),
		[]byte(`{
				"id": "a",
				"login": "karpovIvan@mail.ru",
				"name": "IvAn",
				"password": "QWERTY1234",
				"surname": "Karpov"
			}`),
		[]byte(`{
				"id": 0,
				"login": "n@mail.ru",
				"name": "Nikita",
				"password": "qwerty1234",
				"surname": "Nosov"
			}`),
	}

	for _, body := range arrBody {
		err := registerUser(t, userHandler, body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	expectedStatus := []int{200, 400, 401}
	for i, body := range arrBadStatusBody {
		r := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		userHandler.Login(w, r)
		fmt.Println(w.Code)
		if w.Code != expectedStatus[i] {
			assert.Equal(t, expectedStatus[i], w.Code)
			t.Error("status is not ok")
		}
	}
}

func TestLogoutUser(t *testing.T) {
	t.Parallel()

	var (
		userRepository = user.NewEmptyInMemoryUserRepository()
		userHandler    = &UserHandler{
			UserRepository: userRepository,
		}
	)

	for _, body := range arrBody {
		err := registerUser(t, userHandler, body)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var cookies []string
	for _, body := range arrBody {
		cookie, err := loginUser(t, userHandler, body)
		cookies = append(cookies, cookie)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for _, c := range cookies {
		r := httptest.NewRequest("POST", "/logout", nil)
		fmt.Println("c: ", c)
		cookie := &http.Cookie{
			Name:    "session_id",
			Value:   c,
			Expires: time.Now().Add(90 * 24 * time.Hour),
			Path:    "/",
		}
		r.AddCookie(cookie)
		w := httptest.NewRecorder()

		userHandler.Logout(w, r)
		fmt.Println(w.Code)
		if w.Code != http.StatusOK {
			assert.Equal(t, http.StatusOK, w.Code)
			t.Error("status is not ok")
		}
	}
}
