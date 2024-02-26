package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mail/pkg/email"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// EmailHandler represents the handler for email operations.
type EmailHandler struct {
	EmailRepository email.EmailRepository
}

// List displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Produce json
// @Success 200 {array} map[string]uint64
// @Router /emails [get]
func (h *EmailHandler) List(w http.ResponseWriter, r *http.Request) {
	/* emails, err := h.EmailRepository.GetAll()
	if err != nil {
		http.Error(w, `DB error`, http.StatusInternalServerError)
		return
	}

	respJSON, err := json.Marshal(emails)*/
	_, err := SM.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	respJSON, err := json.Marshal(email.FakeEmails)
	if err != nil {
		http.Error(w, `JSON encoding error`, http.StatusInternalServerError)
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
// @Router /email/{id} [get]
func (h *EmailHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	_, err := SM.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	email, err := h.EmailRepository.GetByID(id)
	if err != nil {
		http.Error(w, `DB error`, http.StatusInternalServerError)
		return
	}
	if email == nil {
		http.Error(w, `No email found`, http.StatusNotFound)
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
// @Router /email/add [post]
func (h *EmailHandler) Add(w http.ResponseWriter, r *http.Request) {
	_, err := SM.Check(r)
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
		http.Error(w, `Bad JSON`, http.StatusBadRequest)
		return
	}

	id, err := h.EmailRepository.Add(&newEmail)
	if err != nil {
		http.Error(w, `DB error`, http.StatusInternalServerError)
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
// @Router /email/update/{id} [put]
func (h *EmailHandler) Update(w http.ResponseWriter, r *http.Request) {
	_, err := SM.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	var updatedEmail email.Email
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	err = json.NewDecoder(r.Body).Decode(&updatedEmail)
	if err != nil {
		http.Error(w, `Bad JSON`, http.StatusBadRequest)
		return
	}
	updatedEmail.ID = id

	ok, err := h.EmailRepository.Update(&updatedEmail)
	if err != nil {
		http.Error(w, `DB error`, http.StatusInternalServerError)
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
// @Router /email/delete/{id} [delete]
func (h *EmailHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := SM.Check(r)
	if err != nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, `Bad id`, http.StatusBadRequest)
		return
	}

	ok, err := h.EmailRepository.Delete(id)
	if err != nil {
		http.Error(w, `DB error`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	respJSON, _ := json.Marshal(map[string]bool{"success": ok})
	w.Write(respJSON)
}
