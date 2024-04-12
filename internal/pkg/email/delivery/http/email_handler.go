package http

import (
	"encoding/json"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	"regexp"
	"strconv"

	converters "mail/internal/models/delivery_converters"

	emailApi "mail/internal/models/delivery_models"
	response "mail/internal/models/response"
	emailUsecase "mail/internal/pkg/email/interface"
	domainSession "mail/internal/pkg/session/interface"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var (
	EHandler                        = &EmailHandler{}
	requestIDContextKey interface{} = "requestid"
)

// EmailHandler represents the handler for email operations.
type EmailHandler struct {
	EmailUseCase emailUsecase.EmailUseCase
	Sessions     domainSession.SessionsManager
}

func sanitizeString(str string) string {
	p := bluemonday.UGCPolicy()
	p.AllowElements("b", "i", "a", "strong", "em", "p", "br", "span", "ul", "ol", "li", "h1", "h2", "h3", "div")
	return p.Sanitize(str)
}

// Incoming displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 404 {object} response.Response "DB error"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/emails/incoming [get]
func (h *EmailHandler) Incoming(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	login, err := h.Sessions.GetLoginBySession(r, requestID)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = h.Sessions.CheckLogin(login, requestID, r)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	emails, err := h.EmailUseCase.GetAllEmailsIncoming(login, requestID, 0, 0)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	emailsApi := make([]*emailApi.Email, 0, len(emails))
	for _, email := range emails {
		emailsApi = append(emailsApi, converters.EmailConvertCoreInApi(*email))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// Sent displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 404 {object} response.Response "DB error"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/emails/sent [get]
func (h *EmailHandler) Sent(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	login, err := h.Sessions.GetLoginBySession(r, requestID)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = h.Sessions.CheckLogin(login, requestID, r)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	emails, err := h.EmailUseCase.GetAllEmailsSent(login, requestID, 0, 0)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	emailsApi := make([]*emailApi.Email, 0, len(emails))
	for _, email := range emails {
		emailsApi = append(emailsApi, converters.EmailConvertCoreInApi(*email))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// GetByID returns an email message by its ID.
// @Summary Get an email message by ID
// @Description Get an email message by its unique identifier
// @Tags emails
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Email message data"
// @Failure 400 {object} response.Response "Bad id in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 404 {object} response.Response "Email not found"
// @Router /api/v1/email/{id} [get]
func (h *EmailHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	login, err := h.Sessions.GetLoginBySession(r, requestID)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = h.Sessions.CheckLogin(login, requestID, r)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	email, err := h.EmailUseCase.GetEmailByID(id, login, requestID)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, "Email not found")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})
}

// Send adds a new email message.
// @Summary Send a new email message
// @Description Send a new email message to the system
// @Tags emails
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param email body response.EmailSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "ID of the send email message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email message"
// @Router /api/v1/email/send [post]
func (h *EmailHandler) Send(w http.ResponseWriter, r *http.Request) {
	var newEmail emailApi.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newEmail.Topic = sanitizeString(newEmail.Topic)
	newEmail.Text = sanitizeString(newEmail.Text)
	newEmail.PhotoID = sanitizeString(newEmail.PhotoID)
	newEmail.SenderEmail = sanitizeString(newEmail.SenderEmail)
	newEmail.RecipientEmail = sanitizeString(newEmail.RecipientEmail)

	sender := newEmail.SenderEmail
	recipient := newEmail.RecipientEmail
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	switch {
	case isValidMailhubFormat(sender) && isValidMailhubFormat(recipient):
		err = h.Sessions.CheckLogin(sender, requestID, r)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad sender login")
			return
		}

		err = h.EmailUseCase.CheckRecipientEmail(recipient, requestID)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad login")
			return
		}

		email_id, email, err := h.EmailUseCase.CreateEmail(converters.EmailConvertApiInCore(newEmail), requestID)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		err = h.EmailUseCase.CreateProfileEmail(email_id, sender, recipient, requestID)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})
		return
	case isValidMailhubFormat(sender) == true && isValidMailhubFormat(recipient) == false:
		/*email_id, email, err := h.EmailUseCase.CreateEmail(converters.EmailConvertApiInCore(newEmail))
		if err != nil {
			delivery.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})*/
		response.HandleSuccess(w, http.StatusBadRequest, "An error occurred in the recipient's domain. You cannot send messages to other email services. Make sure that the recipient's domain ends with @mailhub.su")
		return
	case isValidMailhubFormat(sender) == false && isValidMailhubFormat(recipient) == true:
		err = h.EmailUseCase.CheckRecipientEmail(recipient, requestID)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad login")
			return
		}

		email_id, email, err := h.EmailUseCase.CreateEmail(converters.EmailConvertApiInCore(newEmail), requestID)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		err = h.EmailUseCase.CreateProfileEmail(email_id, sender, recipient, requestID)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})
		return
	}

	/*email, err := h.EmailUseCase.CreateEmail(converters.EmailConvertApiInCore(newEmail))
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})*/
}

// Update updates an existing email message.
// @Summary Update an email message
// @Description Update an existing email message based on its identifier
// @Tags emails
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path integer true "ID of the email message"
// @Param email body response.EmailSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "Update success status"
// @Failure 400 {object} response.Response "Bad id or Bad JSON"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to update email message"
// @Router /api/v1/email/update/{id} [put]
func (h *EmailHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	var updatedEmail emailApi.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = json.NewDecoder(r.Body).Decode(&updatedEmail)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	updatedEmail.Topic = sanitizeString(updatedEmail.Topic)
	updatedEmail.Text = sanitizeString(updatedEmail.Text)
	updatedEmail.PhotoID = sanitizeString(updatedEmail.PhotoID)
	updatedEmail.RecipientEmail = sanitizeString(updatedEmail.RecipientEmail)
	updatedEmail.SenderEmail = sanitizeString(updatedEmail.SenderEmail)

	updatedEmail.ID = id
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	err1 := h.Sessions.CheckLogin(updatedEmail.SenderEmail, requestID, r)
	err2 := h.Sessions.CheckLogin(updatedEmail.RecipientEmail, requestID, r)
	if err1 != nil && err2 != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	ok, err = h.EmailUseCase.UpdateEmail(converters.EmailConvertApiInCore(updatedEmail), requestID)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to update email message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": ok})
}

// Delete deletes an email message.
// @Summary Delete an email message
// @Description Delete an email message based on its identifier
// @Tags emails
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Deletion success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete email message"
// @Router /api/v1/email/delete/{id} [delete]
func (h *EmailHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	login, err := h.Sessions.GetLoginBySession(r, requestID)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = h.Sessions.CheckLogin(login, requestID, r)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	ok, err = h.EmailUseCase.DeleteEmail(id, login, requestID)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete email message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": ok})
}

func (h *EmailHandler) SendFromAnotherDomain(w http.ResponseWriter, r *http.Request) {
	h.Send(w, r)
}

func isValidMailhubFormat(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@mailhub\.su$`)
	return emailRegex.MatchString(email)
}
