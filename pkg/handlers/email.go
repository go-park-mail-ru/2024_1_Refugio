package handlers

import (
	"encoding/json"
	"fmt"
	"mail/pkg/session"
	"net/http"
	"strconv"

	"mail/pkg/email"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// EmailHandler represents the handler for email operations.
type EmailHandler struct {
	EmailRepository email.EmailRepository
	Sessions        *session.SessionsManager
}

// List displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Produce json
// @Success 200 {object} Response "List of all email messages"
// @Failure 401 {object} Response "Not Authorized"
// @Failure 404 {object} Response "DB error"
// @Failure 500 {object} Response "JSON encoding error"
// @Router /api/v1/emails [get]
func (h *EmailHandler) List(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	emails, err := h.EmailRepository.GetAll()
	if err != nil {
		handleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emails})
}

// GetByID returns an email message by its ID.
// @Summary Get an email message by ID
// @Description Get an email message by its unique identifier
// @Produce json
// @Param id path integer true "ID of the email message"
// @Success 200 {object} Response "Email message data"
// @Failure 400 {object} Response "Bad id in request"
// @Failure 401 {object} Response "Not Authorized"
// @Failure 404 {object} Response "Email not found"
// @Router /api/v1/email/{id} [get]
func (h *EmailHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	email, err := h.EmailRepository.GetByID(id)
	if err != nil {
		handleError(w, http.StatusNotFound, "Email not found")
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"email": email})
}

// Add adds a new email message.
// @Summary Add a new email message
// @Description Add a new email message to the system
// @Accept json
// @Produce json
// @Param email body email.Email true "Email message in JSON format"
// @Success 200 {object} Response "ID of the added email message"
// @Failure 400 {object} Response "Bad JSON in request"
// @Failure 401 {object} Response "Not Authorized"
// @Failure 500 {object} Response "Failed to add email message"
// @Router /api/v1/email/add [post]
func (h *EmailHandler) Add(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	var newEmail email.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	email, err := h.EmailRepository.Add(&newEmail)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to add email message")
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"email": email})
}

// Update updates an existing email message.
// @Summary Update an email message
// @Description Update an existing email message based on its identifier
// @Accept json
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param email body email.Email true "Email message in JSON format"
// @Success 200 {object} Response "Update success status"
// @Failure 400 {object} Response "Bad id or Bad JSON"
// @Failure 401 {object} Response "Not Authorized"
// @Failure 500 {object} Response "Failed to update email message"
// @Router /api/v1/email/update/{id} [put]
func (h *EmailHandler) Update(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	var updatedEmail email.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = json.NewDecoder(r.Body).Decode(&updatedEmail)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}
	updatedEmail.ID = id

	ok, err := h.EmailRepository.Update(&updatedEmail)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to update email message")
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"Success": ok})
}

// Delete deletes an email message.
// @Summary Delete an email message
// @Description Delete an email message based on its identifier
// @Produce json
// @Param id path integer true "ID of the email message"
// @Success 200 {object} Response "Deletion success status"
// @Failure 400 {object} Response "Bad id"
// @Failure 401 {object} Response "Not Authorized"
// @Failure 500 {object} Response "Failed to delete email message"
// @Router /api/v1/email/delete/{id} [delete]
func (h *EmailHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		handleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		handleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	ok, err := h.EmailRepository.Delete(id)
	if err != nil {
		handleError(w, http.StatusInternalServerError, "Failed to delete email message")
		return
	}

	handleSuccess(w, http.StatusOK, map[string]interface{}{"Success": ok})
}
