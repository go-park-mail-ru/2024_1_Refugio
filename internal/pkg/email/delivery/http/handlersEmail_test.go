package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	mock "mail/internal/microservice/email/mock"
	emailCore "mail/internal/microservice/models/domain_models"
	converters "mail/internal/models/delivery_converters"
	emailApi "mail/internal/models/delivery_models"
	"mail/internal/pkg/logger"
	mockSession "mail/internal/pkg/session/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func GetCTX() context.Context {
	requestID := "testID"

	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + "log.txt")
	}
	defer f.Close()

	c := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
	ctx := context.WithValue(c, "requestID", requestID)
	return ctx
}

func TestGelAllIncoming(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		EmailUseCase: mockEmailUseCase,
		Sessions:     mockSessionsManager,
	}

	ctx := GetCTX()

	login := "recipient_test@mailhub.su"

	t.Run("GelAllIncomingSuccess", func(t *testing.T) {
		inemail := &emailCore.Email{
			ID:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		var incominEmails = []*emailCore.Email{inemail}

		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailUseCase.EXPECT().GetAllEmailsIncoming(login, 0, 0, ctx).Return(incominEmails, nil)

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GelAllIncoming Fail in GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllIncoming Fail in CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllIncoming Fail in GetAllEmailsIncoming", func(t *testing.T) {
		inemail := &emailCore.Email{
			ID:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		var incominEmails = []*emailCore.Email{inemail}

		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailUseCase.EXPECT().GetAllEmailsIncoming(login, 0, 0, ctx).Return(incominEmails, fmt.Errorf("GetAllEmailsIncoming"))

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGelAllSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		EmailUseCase: mockEmailUseCase,
		Sessions:     mockSessionsManager,
	}

	ctx := GetCTX()
	login := "sender_test@mailhub.su"

	t.Run("GelAllSentSuccess", func(t *testing.T) {
		inemail := &emailCore.Email{
			ID:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		var incominEmails = []*emailCore.Email{inemail}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailUseCase.EXPECT().GetAllEmailsSent(login, 0, 0, ctx).Return(incominEmails, nil)

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GelAllSent Fail in GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in GetAllEmailsIncoming", func(t *testing.T) {
		inemail := &emailCore.Email{
			ID:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		var incominEmails = []*emailCore.Email{inemail}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailUseCase.EXPECT().GetAllEmailsSent(login, 0, 0, ctx).Return(incominEmails, fmt.Errorf("GetAllEmailsIncoming"))

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		EmailUseCase: mockEmailUseCase,
		Sessions:     mockSessionsManager,
	}
	ctx := GetCTX()
	login := "test@mailhub.su"

	t.Run("GetByID Successs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/email/{id}", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailUseCase.EXPECT().GetEmailByID(uint64(1), login, r.Context()).Return(&emailCore.Email{}, nil)

		emailHandler.GetByID(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetByID Fail GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/email/{id}", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.GetByID(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetByID CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/email/{id}", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(fmt.Errorf("CheckLogin"))

		emailHandler.GetByID(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetByID GetEmailByID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/email/{id}", bytes.NewReader([]byte(``)))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailUseCase.EXPECT().GetEmailByID(uint64(1), login, r.Context()).Return(&emailCore.Email{}, fmt.Errorf("GetEmailByID"))

		emailHandler.GetByID(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		EmailUseCase: mockEmailUseCase,
		Sessions:     mockSessionsManager,
	}
	ctx := GetCTX()

	t.Run("Sender(mailhub.su) Recipient(mailhub.su)", func(t *testing.T) {
		newEmail := emailApi.Email{
			ID:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			PhotoID:        "",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		requestBodyBytes, _ := json.Marshal(newEmail)

		req := httptest.NewRequest("POST", "/api/v1/email/send", bytes.NewReader(requestBodyBytes))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().CheckLogin(newEmail.SenderEmail, r, ctx).Return(nil)
		mockEmailUseCase.EXPECT().CheckRecipientEmail(newEmail.RecipientEmail, ctx).Return(nil)
		mockEmailUseCase.EXPECT().CreateEmail(converters.EmailConvertApiInCore(newEmail), ctx).Return(int64(1), &emailCore.Email{}, nil)
		mockEmailUseCase.EXPECT().CreateProfileEmail(int64(1), newEmail.SenderEmail, newEmail.RecipientEmail, ctx)

		emailHandler.Send(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Sender(mailhub.su) Recipient(another domen)", func(t *testing.T) {
		newEmail := emailApi.Email{
			ID:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			PhotoID:        "",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mail.ru",
		}

		requestBodyBytes, _ := json.Marshal(newEmail)

		req := httptest.NewRequest("POST", "/api/v1/email/send", bytes.NewReader(requestBodyBytes))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		emailHandler.Send(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Sender(another domen) Recipient(mailhub.su)", func(t *testing.T) {
		newEmail := emailApi.Email{
			ID:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			PhotoID:        "",
			SenderEmail:    "sender_test@mail.ru",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		requestBodyBytes, _ := json.Marshal(newEmail)

		req := httptest.NewRequest("POST", "/api/v1/email/send", bytes.NewReader(requestBodyBytes))
		r := req.WithContext(ctx)
		w := httptest.NewRecorder()

		mockEmailUseCase.EXPECT().CheckRecipientEmail(newEmail.RecipientEmail, ctx).Return(nil)
		mockEmailUseCase.EXPECT().CreateEmail(converters.EmailConvertApiInCore(newEmail), ctx).Return(int64(1), &emailCore.Email{}, nil)
		mockEmailUseCase.EXPECT().CreateProfileEmail(int64(1), newEmail.SenderEmail, newEmail.RecipientEmail, ctx)

		emailHandler.Send(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		EmailUseCase: mockEmailUseCase,
		Sessions:     mockSessionsManager,
	}
	ctx := GetCTX()

	newEmail := emailApi.Email{
		ID:             uint64(1),
		Topic:          "Hello",
		Text:           "Hello Sergey",
		PhotoID:        "",
		SenderEmail:    "sender_test@mailhub.su",
		RecipientEmail: "recipient_test@mailhub.su",
	}
	requestBodyBytes, _ := json.Marshal(newEmail)

	t.Run("Update Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/update/{id}", bytes.NewReader(requestBodyBytes))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().CheckLogin(newEmail.SenderEmail, r, r.Context()).Return(nil)
		mockSessionsManager.EXPECT().CheckLogin(newEmail.RecipientEmail, r, r.Context()).Return(nil)
		mockEmailUseCase.EXPECT().UpdateEmail(converters.EmailConvertApiInCore(newEmail), r.Context()).Return(true, nil)

		emailHandler.Update(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Fail", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/update/{id}", bytes.NewReader(requestBodyBytes))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().CheckLogin(newEmail.SenderEmail, r, r.Context()).Return(fmt.Errorf("CheckLogin"))
		mockSessionsManager.EXPECT().CheckLogin(newEmail.RecipientEmail, r, r.Context()).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Update(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update UpdateEmail", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/update/{id}", bytes.NewReader(requestBodyBytes))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().CheckLogin(newEmail.SenderEmail, r, r.Context()).Return(nil)
		mockSessionsManager.EXPECT().CheckLogin(newEmail.RecipientEmail, r, r.Context()).Return(nil)
		mockEmailUseCase.EXPECT().UpdateEmail(converters.EmailConvertApiInCore(newEmail), r.Context()).Return(false, fmt.Errorf("UpdateEmail"))

		emailHandler.Update(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailUseCase := mock.NewMockEmailUseCase(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		EmailUseCase: mockEmailUseCase,
		Sessions:     mockSessionsManager,
	}

	ctx := GetCTX()
	login := "test@mailhub.su"

	t.Run("Delete Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/delete/{id}", bytes.NewReader([]byte("")))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailUseCase.EXPECT().DeleteEmail(uint64(1), login, r.Context()).Return(true, nil)

		emailHandler.Delete(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Fail GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/delete/{id}", bytes.NewReader([]byte("")))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.Delete(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Delete Fail CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/delete/{id}", bytes.NewReader([]byte("")))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Delete(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Delete Fail DeleteEmail", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/delete/{id}", bytes.NewReader([]byte("")))
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)
		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailUseCase.EXPECT().DeleteEmail(uint64(1), login, r.Context()).Return(false, fmt.Errorf("DeleteEmail"))

		emailHandler.Delete(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestIsValidMailhubFormat(t *testing.T) {
	t.Run("IsValidMailhubFormat Success", func(t *testing.T) {
		login := "test@mailhub.su"
		assert.True(t, isValidMailhubFormat(login))
	})

	t.Run("IsValidMailhubFormat Fail", func(t *testing.T) {
		login := "test@mail.su"
		assert.False(t, isValidMailhubFormat(login))
	})
}

func TestSanitizeString(t *testing.T) {
	t.Run("IsValidMailhubFormat Ease", func(t *testing.T) {
		s := "OK"
		answer := sanitizeString(s)
		assert.Equal(t, s, answer)
	})

	t.Run("IsValidMailhubFormat Span", func(t *testing.T) {
		s := "</span>OK<span>"
		answer := sanitizeString(s)
		assert.Equal(t, s, answer)
	})

	t.Run("IsValidMailhubFormat Script", func(t *testing.T) {
		s := "</script>OK<script>"
		answer := sanitizeString(s)
		assert.Equal(t, "OK", answer)
	})
}
