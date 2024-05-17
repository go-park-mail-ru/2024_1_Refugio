package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/striphtmltags"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"google.golang.org/api/gmail/v1"
	apiModels "mail/internal/models/delivery_models"
	"mail/internal/models/response"
	"mail/internal/pkg/utils/validators"
	"net/http"
)

// AddDraft adds a new draft email message.
// @Summary AddDraft a new draft email message
// @Description AddDraft a new draft email message to the system
// @Tags drafts-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param email body response.EmailOtherSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "ID of the send email message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email message"
// @Router /api/v1/gmail/draft/adddraft [post]
func (g *GMailEmailHandler) AddDraft(w http.ResponseWriter, r *http.Request) {
	var newEmail apiModels.OtherEmail
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = g.Sessions.CheckLogin(newEmail.SenderEmail, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	strRaw := fmt.Sprintf("From: %v\r\n", newEmail.SenderEmail)
	if !validators.IsEmpty(newEmail.RecipientEmail) {
		strRaw += fmt.Sprintf("To: %v\r\n", newEmail.RecipientEmail)
	}
	strRaw += fmt.Sprintf("Subject: %v\r\n\r\n%v", newEmail.Topic, newEmail.Text)

	input := base64.RawStdEncoding.EncodeToString([]byte(strRaw))

	draft := &gmail.Draft{
		Message: &gmail.Message{
			Raw: input,
		},
	}

	_, err = srv.Users.Drafts.Create("me", draft).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error create the draft")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newEmail})
}

// SendDraft adds a new sent email message.
// @Summary SendDraft a new sent email message
// @Description AddDraft a new sent email message to the system
// @Tags drafts-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param email body response.EmailOtherSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "ID of the send email message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email message"
// @Router /api/v1/gmail/draft/sendDraft [post]
func (g *GMailEmailHandler) SendDraft(w http.ResponseWriter, r *http.Request) {
	var newEmail apiModels.OtherEmail
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = g.Sessions.CheckLogin(newEmail.SenderEmail, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	if validators.IsEmpty(newEmail.RecipientEmail) {
		response.HandleError(w, http.StatusBadRequest, "Empty login recipient")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	input := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("From: %v\r\nTo: %v\r\nSubject: %v\r\n\r\n%v", newEmail.SenderEmail, newEmail.RecipientEmail, newEmail.Topic, newEmail.Text)))
	draft := &gmail.Draft{
		Id: newEmail.ID,
		Message: &gmail.Message{
			Raw: input,
		},
	}

	_, err = srv.Users.Drafts.Send("me", draft).Do()
	if err != nil {
		fmt.Println(err)
		response.HandleError(w, http.StatusInternalServerError, "Error send the draft")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newEmail})
}

// GetByIdDraft returns an draft message by its ID.
// @Summary Get an draft message by ID
// @Description Get an email message by its unique identifier
// @Tags drafts-gmail
// @Produce json
// @Param id path string true "ID of the draft message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Email message data"
// @Failure 400 {object} response.Response "Bad id in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 404 {object} response.Response "Email not found"
// @Router /api/v1/gmail/draft/{id} [get]
func (g *GMailEmailHandler) GetByIdDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	draftID, ok := vars["id"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user session")
		return
	}

	if !validators.IsValidEmailFormatGmail(login) {
		response.HandleError(w, http.StatusBadRequest, "Login must end with @gmail.com")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	draft, err := srv.Users.Drafts.Get("me", draftID).Format("full").Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error receiving messages")
		return
	}

	email := CreateEmailStructDraft(draft)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": email})
}

// UpdateDraft update a draft message.
// @Summary SendDraft a update draft message
// @Description AddDraft a nupdate draft message to the system
// @Tags drafts-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path string true "ID of the folder message"
// @Param email body response.EmailOtherSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "ID of the send email message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email message"
// @Router /api/v1/gmail/draft/update/{id} [put]
func (g *GMailEmailHandler) UpdateDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	draftID, ok := vars["id"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	var newEmail apiModels.OtherEmail
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = g.Sessions.CheckLogin(newEmail.SenderEmail, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	strRaw := fmt.Sprintf("From: %v\r\n", newEmail.SenderEmail)
	if !validators.IsEmpty(newEmail.RecipientEmail) {
		strRaw += fmt.Sprintf("To: %v\r\n", newEmail.RecipientEmail)
	}
	strRaw += fmt.Sprintf("Subject: %v\r\n\r\n%v", newEmail.Topic, newEmail.Text)

	input := base64.RawStdEncoding.EncodeToString([]byte(strRaw))

	draft := &gmail.Draft{
		Id: draftID,
		Message: &gmail.Message{
			Raw: input,
		},
	}

	_, err = srv.Users.Drafts.Update("me", draftID, draft).Do()
	if err != nil {
		fmt.Println(err)
		response.HandleError(w, http.StatusInternalServerError, "Error update the draft")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newEmail})
}

// GetDrafts displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags drafts-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/gmail/drafts [get]
func (g *GMailEmailHandler) GetDrafts(w http.ResponseWriter, r *http.Request) {
	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user session")
		return
	}
	if !validators.IsValidEmailFormatGmail(login) {
		response.HandleError(w, http.StatusBadRequest, "Login must end with @gmail.com")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	req, err := srv.Users.Drafts.List("me").Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error receiving list messages")
		return
	}

	emailsApi := make([]*apiModels.OtherEmail, len(req.Drafts))
	for i, d := range req.Drafts {
		dr, err := srv.Users.Drafts.Get("me", d.Id).Format("full").Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error receiving messages")
			return
		}
		email := CreateEmailStructDraft(dr)
		emailsApi[i] = email
	}

	for i, _ := range emailsApi {
		emailsApi[i].DraftStatus = true
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// DeleteDraft deletes an draft message.
// @Summary Delete an draft message
// @Description Delete an draft message based on its identifier
// @Tags drafts-gmail
// @Produce json
// @Param id path string true "ID of the draft message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Deletion success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete email message"
// @Router /api/v1/gmail/draft/delete/{id} [delete]
func (g *GMailEmailHandler) DeleteDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	draftID, ok := vars["id"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user session")
		return
	}

	if !validators.IsValidEmailFormatGmail(login) {
		response.HandleError(w, http.StatusBadRequest, "Login must end with @gmail.com")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	err = srv.Users.Drafts.Delete("me", draftID).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error delete draft")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": true})
}

func CreateEmailStructDraft(msg *gmail.Draft) *apiModels.OtherEmail {
	email := &apiModels.OtherEmail{}
	email.ID = msg.Id

	fmt.Println(msg.Message.Payload.MimeType)
	if msg.Message.Payload.MimeType == "text/plain" {
		email = ParserMessageHeadresDraft(email, msg)
		data, err := base64.URLEncoding.DecodeString(msg.Message.Payload.Body.Data)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		email.Text = striphtmltags.StripTags(string(data))
	} else if msg.Message.Payload.MimeType == "text/html" {
		data, err := base64.URLEncoding.DecodeString(msg.Message.Payload.Body.Data)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		email.Text = string(data)
		email = ParserMessageHeadresDraft(email, msg)
	} else if len(msg.Message.Payload.Parts) != 0 {
		for _, part := range msg.Message.Payload.Parts {
			if part.MimeType == "text/html" {
				data, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					fmt.Println("Error: ", err)
				}
				email.Text = string(data)
				email = ParserMessageHeadresDraft(email, msg)
			}
		}
	}

	return email
}

func ParserMessageHeadresDraft(email *apiModels.OtherEmail, msg *gmail.Draft) *apiModels.OtherEmail {
	for _, mes := range msg.Message.Payload.Headers {
		if mes.Name == "To" {
			email.RecipientEmail = mes.Value
		}
		if mes.Name == "From" {
			email.SenderEmail = mes.Value
		}
		if mes.Name == "Subject" {
			email.Topic = striphtmltags.StripTags(mes.Value)
		}
	}
	return email
}
