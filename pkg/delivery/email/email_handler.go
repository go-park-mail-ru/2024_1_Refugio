package email

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"mail/pkg/delivery"
	"mail/pkg/delivery/converters"

	emailApi "mail/pkg/delivery/models"
	"mail/pkg/delivery/session"
	emailUsecase "mail/pkg/domain/usecase"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var (
	EHandler = &EmailHandler{}
)

// EmailHandler represents the handler for email operations.
type EmailHandler struct {
	EmailUseCase emailUsecase.EmailUseCase
	Sessions     *session.SessionsManager
}

func InitializationEmailHandler(emailHandler *EmailHandler) {
	EHandler = emailHandler
}

// List displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails
// @Produce json
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "List of all email messages"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Failure 404 {object} delivery.Response "DB error"
// @Failure 500 {object} delivery.Response "JSON encoding error"
// @Router /api/v1/auth/emails [get]
func (h *EmailHandler) List(w http.ResponseWriter, r *http.Request) {
	emails, err := h.EmailUseCase.GetAllEmails(0, 0)
	if err != nil {
		delivery.HandleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	emailsApi := make([]*emailApi.Email, 0, len(emails))
	for _, email := range emails {
		emailsApi = append(emailsApi, converters.EmailConvertCoreInApi(*email))
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// GetByID returns an email message by its ID.
// @Summary Get an email message by ID
// @Description Get an email message by its unique identifier
// @Tags emails
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "Email message data"
// @Failure 400 {object} delivery.Response "Bad id in request"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Failure 404 {object} delivery.Response "Email not found"
// @Router /api/v1/auth/email/{id} [get]
func (h *EmailHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	email, err := h.EmailUseCase.GetEmailByID(id)
	if err != nil {
		delivery.HandleError(w, http.StatusNotFound, "Email not found")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})
}

// Add adds a new email message.
// @Summary Add a new email message
// @Description Add a new email message to the system
// @Tags emails
// @Accept json
// @Produce json
// @Param email body delivery.EmailSwag true "Email message in JSON format"
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "ID of the added email message"
// @Failure 400 {object} delivery.Response "Bad JSON in request"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Failure 500 {object} delivery.Response "Failed to add email message"
// @Router /api/v1/auth/email/add [post]
func (h *EmailHandler) Add(w http.ResponseWriter, r *http.Request) {
	var newEmail emailApi.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	email, err := h.EmailUseCase.CreateEmail(converters.EmailConvertApiInCore(newEmail))
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})
}

// Update updates an existing email message.
// @Summary Update an email message
// @Description Update an existing email message based on its identifier
// @Tags emails
// @Accept json
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param email body delivery.EmailSwag true "Email message in JSON format"
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "Update success status"
// @Failure 400 {object} delivery.Response "Bad id or Bad JSON"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Failure 500 {object} delivery.Response "Failed to update email message"
// @Router /api/v1/auth/email/update/{id} [put]
func (h *EmailHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	var updatedEmail emailApi.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = json.NewDecoder(r.Body).Decode(&updatedEmail)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}
	updatedEmail.ID = id

	ok, err := h.EmailUseCase.UpdateEmail(converters.EmailConvertApiInCore(updatedEmail))
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to update email message")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": ok})
}

// Delete deletes an email message.
// @Summary Delete an email message
// @Description Delete an email message based on its identifier
// @Tags emails
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param X-CSRF-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "Deletion success status"
// @Failure 400 {object} delivery.Response "Bad id"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Failure 500 {object} delivery.Response "Failed to delete email message"
// @Router /api/v1/auth/email/delete/{id} [delete]
func (h *EmailHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	ok, err := h.EmailUseCase.DeleteEmail(id)
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to delete email message")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": ok})
}

/* func (h *EmailHandler) Send(w http.ResponseWriter, r *http.Request) {
	// ivan@mailhub.ru -> sergey@mailhub.ru
	if (email1 == @mailhub.ru && email2 == @mailhub.ru) {
		h.EmailUseCase.CreateEmail()
	}

	// ivan@mailhub.ru -> sergey@yandex.ru
	if (email1 == @mailhub.ru && email2 == @yandex.ru) {
		// smtp
		h.EmailUseCase.CreateEmail()
	}

	// ivan@yandex.ru -> sergey@mailhub.ru
	if (email1 == @yandex.ru && email2 == @mailhub.ru) {
		h.EmailUseCase.CreateEmail()
	}
} */
