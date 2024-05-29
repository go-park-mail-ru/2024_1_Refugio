package http

import (
	"encoding/base64"
	"errors"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/denisbrodbeck/striphtmltags"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"mail/internal/models/response"
	"mail/internal/pkg/utils/validators"

	apiModels "mail/internal/models/delivery_models"
	gmailAuth "mail/internal/pkg/gmail/gmail_auth/delivery/http"
	domainSession "mail/internal/pkg/session/interface"
)

// GMailEmailHandler handles user-related HTTP requests.
type GMailEmailHandler struct {
	Sessions     domainSession.SessionsManager
	GMailService *gmail.Service
}

func sanitizeString(str string) string {
	p := bluemonday.UGCPolicy()
	p.AllowElements("b", "i", "a", "strong", "em", "p", "br", "span", "ul", "ol", "li", "h1", "h2", "h3", "div")
	return p.Sanitize(str)
}

// GetIncoming displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/gmail/emails/incoming [get]
func (g *GMailEmailHandler) GetIncoming(w http.ResponseWriter, r *http.Request) {
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

	req, err := srv.Users.Messages.List("me").Q("label:inbox").Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error receiving list messages")
		return
	}

	p := bluemonday.StripTagsPolicy()
	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error receiving messages")
			return
		}

		email, err := CreateEmailStruct(msg)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error decoding body data")
			return
		}

		text := p.Sanitize(email.Text)
		text = strings.ReplaceAll(text, "\n", "")
		fields := strings.Fields(text)
		email.Text = strings.Join(fields, " ")
		for _, l := range msg.LabelIds {
			label, err := srv.Users.Labels.Get("me", l).Do()
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed get label")
				return
			}
			if label.Name == "UNREAD" {
				email.ReadStatus = false
				break
			} else {
				email.ReadStatus = true
			}
		}
		emailsApi[i] = email
	}
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// GetSent displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/gmail/emails/sent [get]
func (g *GMailEmailHandler) GetSent(w http.ResponseWriter, r *http.Request) {
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

	req, err := srv.Users.Messages.List("me").Q("label:sent").Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error receiving list messages")
		return
	}

	p := bluemonday.StripTagsPolicy()
	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error receiving messages")
			return
		}
		email, err := CreateEmailStruct(msg)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error decoding body data")
			return
		}
		text := p.Sanitize(email.Text)
		text = strings.ReplaceAll(text, "\n", "")
		fields := strings.Fields(text)
		email.Text = strings.Join(fields, " ")
		for _, l := range msg.LabelIds {
			label, err := srv.Users.Labels.Get("me", l).Do()
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed get label")
				return
			}
			if label.Name == "UNREAD" {
				email.ReadStatus = false
				break
			} else {
				email.ReadStatus = true
			}
		}
		emailsApi[i] = email
	}
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// GetSpam displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/gmail/emails/spam [get]
func (g *GMailEmailHandler) GetSpam(w http.ResponseWriter, r *http.Request) {
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

	req, err := srv.Users.Messages.List("me").Q("label:spam").Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error receiving list messages")
		return
	}

	p := bluemonday.StripTagsPolicy()
	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error receiving messages")
			return
		}
		email, err := CreateEmailStruct(msg)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error decoding body data")
			return
		}
		text := p.Sanitize(email.Text)
		text = strings.ReplaceAll(text, "\n", "")
		fields := strings.Fields(text)
		email.Text = strings.Join(fields, " ")
		for _, l := range msg.LabelIds {
			label, err := srv.Users.Labels.Get("me", l).Do()
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed get label")
				return
			}
			if label.Name == "UNREAD" {
				email.ReadStatus = false
				break
			} else {
				email.ReadStatus = true
			}
		}
		email.SpamStatus = true
		emailsApi[i] = email
	}

	for i := range emailsApi {
		emailsApi[i].SpamStatus = true
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// GetById returns an email message by its ID.
// @Summary Get an email message by ID
// @Description Get an email message by its unique identifier
// @Tags emails-gmail
// @Produce json
// @Param id path string true "ID of the email message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Email message data"
// @Failure 400 {object} response.Response "Bad id in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 404 {object} response.Response "Email not found"
// @Router /api/v1/gmail/email/{id} [get]
func (g *GMailEmailHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID, ok := vars["id"]
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

	msg, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error receiving messages")
		return
	}

	for _, header := range msg.Payload.Headers {
		fmt.Println("Name: ", header.Name, "   Value: ", header.Value)
		if header.Name == "References" {
			references := strings.Split(header.Value, " ")
			parentMessageId := strings.TrimPrefix(references[0], "message-id:")
			fmt.Println("Идентификатор родительского сообщения:", parentMessageId)
		}
	}

	email, err := CreateEmailStruct(msg)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error decoding body data")
		return
	}

	read := true
	spam := true
	for _, l := range msg.LabelIds {
		label, err := srv.Users.Labels.Get("me", l).Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed get label")
			return
		}
		if read {
			if label.Name == "UNREAD" {
				read = false
			}
		}
		if spam {
			if label.Name == "SPAM" {
				spam = false
			}
		}
	}
	if read {
		email.ReadStatus = true
	}
	if !spam {
		email.SpamStatus = true
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": email})
}

// Send adds a new email message.
// @Summary Send a new email message
// @Description Send a new email message to the system
// @Tags emails-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param email body response.EmailOtherSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "ID of the send email message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email message"
// @Router /api/v1/gmail/email/send [post]
func (g *GMailEmailHandler) Send(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newEmail apiModels.OtherEmail
	if err := newEmail.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newEmail.Topic = sanitizeString(newEmail.Topic)
	newEmail.Text = sanitizeString(newEmail.Text)
	newEmail.SenderEmail = sanitizeString(newEmail.SenderEmail)
	newEmail.RecipientEmail = sanitizeString(newEmail.RecipientEmail)

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
		response.HandleError(w, http.StatusBadRequest, "Recipient is empty")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	input := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("From: %v\r\nTo: %v\r\nSubject: %v\r\n\r\n%v", newEmail.SenderEmail, newEmail.RecipientEmail, newEmail.Topic, newEmail.Text)))

	message := &gmail.Message{
		Raw: input,
	}

	_, err = srv.Users.Messages.Send("me", message).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error sending the message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newEmail})
}

// Delete deletes an email message.
// @Summary Delete an email message
// @Description Delete an email message based on its identifier
// @Tags emails-gmail
// @Produce json
// @Param id path string true "ID of the email message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Deletion success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete email message"
// @Router /api/v1/gmail/email/delete/{id} [delete]
func (g *GMailEmailHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID, ok := vars["id"]
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

	err = srv.Users.Messages.Delete("me", messageID).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error deleting a message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": true})
}

// Update an email message.
// @Summary Update an email draft message
// @Description Update an update email message to the system
// @Tags emails-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path string true "ID of the email message"
// @Param email body response.EmailOtherSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "ID of the update email message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to update email message"
// @Router /api/v1/gmail/email/update/{id} [put]
func (g *GMailEmailHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageID, ok := vars["id"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newEmail apiModels.OtherEmail
	if err := newEmail.UnmarshalJSON(body); err != nil {
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

	var addModify []string
	var removeModify []string
	if newEmail.SpamStatus {
		addModify = append(addModify, "SPAM")
	} else {
		addModify = append(addModify, "INBOX")
		removeModify = append(removeModify, "SPAM")
	}

	if !newEmail.ReadStatus {
		addModify = append(addModify, "UNREAD")
	} else {
		removeModify = append(removeModify, "UNREAD")
	}

	if newEmail.Flag {
		addModify = append(addModify, "IMPORTANT")
	} else {
		removeModify = append(removeModify, "IMPORTANT")
	}

	modifyRequest := &gmail.ModifyMessageRequest{
		RemoveLabelIds: removeModify,
		AddLabelIds:    addModify,
	}
	_, err = srv.Users.Messages.Modify("me", messageID, modifyRequest).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error update the draft")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newEmail})
}

func CreateEmailStruct(msg *gmail.Message) (*apiModels.OtherEmail, error) {
	email := &apiModels.OtherEmail{}
	email.ID = msg.Id
	email.DateOfDispatch = time.Unix(msg.InternalDate/1000, 0)

	if msg.Payload.MimeType == "text/plain" {
		email = ParserMessageHeaders(email, msg)
		data, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			return nil, fmt.Errorf("error decoding body data")
		}
		email.Text = striphtmltags.StripTags(string(data))
	} else if msg.Payload.MimeType == "text/html" {
		data, err := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		if err != nil {
			return nil, fmt.Errorf("error decoding body data")
		}
		email.Text = string(data)
		email = ParserMessageHeaders(email, msg)
	} else if len(msg.Payload.Parts) != 0 {
		for _, part := range msg.Payload.Parts {
			if part.MimeType == "text/html" {
				data, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					return nil, fmt.Errorf("error decoding body data")
				}
				email.Text = string(data)
				email = ParserMessageHeaders(email, msg)
			}
		}
	}

	return email, nil
}

func ParserMessageHeaders(email *apiModels.OtherEmail, msg *gmail.Message) *apiModels.OtherEmail {
	for _, mes := range msg.Payload.Headers {
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

func GetSRV(login string) (*gmail.Service, error) {
	srv, ok := gmailAuth.MapOAuthCongig[login]
	if !ok {
		return nil, errors.New("unable to retrieve Gmail client")
	}
	return srv, nil
}
