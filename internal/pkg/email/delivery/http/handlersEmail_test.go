package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"mail/internal/pkg/utils/constants"

	email_mock "mail/internal/microservice/email/mock"
	email_proto "mail/internal/microservice/email/proto"
	emailApi "mail/internal/models/delivery_models"
	mockSession "mail/internal/pkg/session/mock"
)

func TestGelAllIncoming(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

	login := "recipient_test@mailhub.su"

	t.Run("GelAllIncomingSuccess", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetAllIncoming(gomock.Any(), gomock.Any()).Return(mockEmails, nil)

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GelAllIncoming Fail in GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllIncoming Fail in CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllIncoming Fail in GetAllEmailsIncoming", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/incoming", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetAllIncoming(gomock.Any(), gomock.Any()).Return(mockEmails, fmt.Errorf("GetAllEmailsIncoming"))

		emailHandler.Incoming(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGelAllSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

	login := "sender_test@mailhub.su"

	t.Run("GelAllSentSuccess", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetAllSent(gomock.Any(), gomock.Any()).Return(mockEmails, nil)

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GelAllSent Fail in GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in GetAllEmailsIncoming", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetAllSent(gomock.Any(), gomock.Any()).Return(mockEmails, fmt.Errorf("GetAllEmailsIncoming"))

		emailHandler.Sent(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

	login := "test@mailhub.su"

	t.Run("GetByID Successs", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/email/{id}", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailServiceClient.EXPECT().GetEmailByID(gomock.Any(), gomock.Any()).Return(&email_proto.Email{}, nil)

		emailHandler.GetByID(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetByID Fail GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/email/{id}", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailServiceClient.EXPECT().GetEmailByID(gomock.Any(), gomock.Any()).Return(&email_proto.Email{}, fmt.Errorf("GetEmailByID"))

		emailHandler.GetByID(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().CheckLogin(newEmail.SenderEmail, r, r.Context()).Return(nil)
		mockSessionsManager.EXPECT().CheckLogin(newEmail.RecipientEmail, r, r.Context()).Return(nil)
		mockEmailServiceClient.EXPECT().UpdateEmail(gomock.Any(), gomock.Any()).Return(&email_proto.StatusEmail{Status: true}, nil)

		emailHandler.Update(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update Fail", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/update/{id}", bytes.NewReader(requestBodyBytes))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().CheckLogin(newEmail.SenderEmail, r, r.Context()).Return(nil)
		mockSessionsManager.EXPECT().CheckLogin(newEmail.RecipientEmail, r, r.Context()).Return(nil)
		mockEmailServiceClient.EXPECT().UpdateEmail(gomock.Any(), gomock.Any()).Return(&email_proto.StatusEmail{Status: false}, fmt.Errorf("UpdateEmail"))

		emailHandler.Update(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

	login := "test@mailhub.su"

	t.Run("Delete Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/delete/{id}", bytes.NewReader([]byte("")))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailServiceClient.EXPECT().DeleteEmail(gomock.Any(), gomock.Any()).Return(&email_proto.StatusEmail{Status: true}, nil)

		emailHandler.Delete(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete Fail GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/email/delete/{id}", bytes.NewReader([]byte("")))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)
		vars := map[string]string{"id": "1"}
		r = mux.SetURLVars(r, vars)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, r.Context()).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, r.Context()).Return(nil)
		mockEmailServiceClient.EXPECT().DeleteEmail(gomock.Any(), gomock.Any()).Return(&email_proto.StatusEmail{Status: false}, fmt.Errorf("DeleteEmail"))

		emailHandler.Delete(w, r)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDraft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

	login := "sender_test@mailhub.su"

	t.Run("GelAllSentSuccess", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetDraftEmails(gomock.Any(), gomock.Any()).Return(mockEmails, nil)

		emailHandler.Draft(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GelAllSent Fail in GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.Draft(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Draft(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in GetAllEmailsIncoming", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetDraftEmails(gomock.Any(), gomock.Any()).Return(mockEmails, fmt.Errorf("GetAllEmailsIncoming"))

		emailHandler.Draft(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSpam(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

	login := "sender_test@mailhub.su"

	t.Run("GelAllSentSuccess", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetSpamEmails(gomock.Any(), gomock.Any()).Return(mockEmails, nil)

		emailHandler.Spam(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GelAllSent Fail in GetLoginBySession", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, fmt.Errorf("GetLoginBySession"))

		emailHandler.Spam(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in CheckLogin", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(fmt.Errorf("CheckLogin"))

		emailHandler.Spam(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GelAllSent Fail in GetAllEmailsIncoming", func(t *testing.T) {
		inemail := &email_proto.Email{
			Id:             uint64(1),
			Topic:          "Hello",
			Text:           "Hello Sergey",
			SenderEmail:    "sender_test@mailhub.su",
			RecipientEmail: "recipient_test@mailhub.su",
		}

		mockEmails := &email_proto.Emails{
			Emails: []*email_proto.Email{inemail},
		}

		req := httptest.NewRequest("GET", "/api/v1/emails/sent", bytes.NewReader([]byte(``)))
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().GetLoginBySession(r, ctx).Return(login, nil)
		mockSessionsManager.EXPECT().CheckLogin(login, r, ctx).Return(nil)
		mockEmailServiceClient.EXPECT().GetSpamEmails(gomock.Any(), gomock.Any()).Return(mockEmails, fmt.Errorf("GetAllEmailsIncoming"))

		emailHandler.Spam(w, r)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

/*
func TestSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailServiceClient := email_mock.NewMockEmailServiceClient(ctrl)
	mockSessionsManager := mockSession.NewMockSessionsManager(ctrl)

	emailHandler := EmailHandler{
		Sessions:           mockSessionsManager,
		EmailServiceClient: mockEmailServiceClient,
	}

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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSessionsManager.EXPECT().CheckLogin(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		mockSessionsManager.EXPECT().CheckLogin(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)                          // Добавим проверку получателя
		mockEmailServiceClient.EXPECT().CheckRecipientEmail(gomock.Any(), gomock.Any()).Return(&email_proto.EmptyEmail{}, nil) // Добавим правильное количество аргументов
		mockEmailServiceClient.EXPECT().CreateEmail(gomock.Any(), gomock.Any()).Return(&email_proto.EmailWithID{Email: nil, Id: uint64(1)}, nil)
		mockEmailServiceClient.EXPECT().CreateProfileEmail(gomock.Any(), gomock.Any()).Return(&email_proto.EmptyEmail{}, nil) // Добавим правильное количество аргументов

		emailHandler.Send(w, req)
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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
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
		ctx := context.WithValue(req.Context(), interface{}(string(constants.RequestIDKey)), "testID")
		r := req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockEmailServiceClient.EXPECT().CheckRecipientEmail(gomock.Any(), gomock.Any()).Return(nil)
		mockEmailServiceClient.EXPECT().CreateEmail(gomock.Any(), gomock.Any()).Return(&email_proto.EmailWithID{Email: nil, Id: uint64(1)}, nil)
		mockEmailServiceClient.EXPECT().CreateProfileEmail(gomock.Any(), gomock.Any()).Return(&email_proto.EmptyEmail{})

		emailHandler.Send(w, r)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
*/

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
