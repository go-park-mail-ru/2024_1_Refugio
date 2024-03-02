package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"mail/pkg/email"
	"mail/pkg/session"

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
// @Success 200 {array} map[string]uint64
// @Failure 401 {string} string "Not Authorized"
// @Failure 404 {string} string "DB error"
// @Failure 500 {string} string "JSON encoding error"
// @Router /emails [get]
func (h *EmailHandler) List(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	dbEmails, err := h.EmailRepository.GetAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusNotFound)
		return
	}
	allEmails := append(dbEmails, email.FakeEmails...)

	respJSON, err := json.Marshal(allEmails)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(respJSON)
}

// GetByID returns an email message by its ID.
// @Summary Get an email message by ID
// @Description Get an email message by its unique identifier
// @Produce json
// @Param id path integer true "ID of the email message"
// @Success 200 {object} map[string]uint64
// @Failure 400 {string} string "Bad id in request"
// @Failure 401 {string} string "Not Authorized"
// @Failure 404 {string} string "DB error"
// @Router /email/{id} [get]
func (h *EmailHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, `Bad id in request`, http.StatusBadRequest)
		return
	}

	email, err := h.EmailRepository.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusNotFound)
		return
	}

	respJSON, _ := json.Marshal(email)
	w.Header().Set("Content-type", "application/json")
	w.Write(respJSON)
}

// Add adds a new email message.
// @Summary Add a new email message
// @Description Add a new email message to the system
// @Accept json
// @Produce json
// @Param email body email.Email true "Email message in JSON format"
// @Success 200 {object} map[string]uint64
// @Failure 401 {string} string "Not Authorized"
// @Failure 400 {string} string "Bad JSON in request"
// @Failure 500 {string} string "DB error"
// @Router /email/add [post]
func (h *EmailHandler) Add(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var newEmail email.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		http.Error(w, `Bad JSON in request`, http.StatusBadRequest)
		return
	}

	id, err := h.EmailRepository.Add(&newEmail)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	respJSON, _ := json.Marshal(map[string]uint64{"id": id})
	w.Header().Set("Content-type", "application/json")
	w.Write(respJSON)
}

// Update updates an existing email message.
// @Summary Update an email message
// @Description Update an existing email message based on its identifier
// @Accept json
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param email body email.Email true "Email message in JSON format"
// @Success 200 {object} map[string]bool
// @Failure 401 {string} string "Not Authorized"
// @Failure 400 {string} string "Bad id or Bad JSON"
// @Failure 500 {string} string "DB error"
// @Router /email/update/{id} [put]
func (h *EmailHandler) Update(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, `Bad id in request`, http.StatusBadRequest)
		return
	}

	var updatedEmail email.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = json.NewDecoder(r.Body).Decode(&updatedEmail)
	if err != nil {
		http.Error(w, `Bad JSON in request`, http.StatusBadRequest)
		return
	}
	updatedEmail.ID = id

	ok, err := h.EmailRepository.Update(&updatedEmail)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]bool{"success": ok})
	w.Write(respJSON)
}

// Delete deletes an email message.
// @Summary Delete an email message
// @Description Delete an email message based on its identifier
// @Produce json
// @Param id path integer true "ID of the email message"
// @Success 200 {object} map[string]bool
// @Failure 400 {string} string "Bad id"
// @Failure 401 {string} string "Not Authorized"
// @Failure 500 {string} string "DB error"
// @Router /email/delete/{id} [delete]
func (h *EmailHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := h.Sessions.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, `Bad id in request`, http.StatusBadRequest)
		return
	}

	ok, err := h.EmailRepository.Delete(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("DB error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]bool{"success": ok})
	w.Write(respJSON)
}
