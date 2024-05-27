package http

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"io"
	"net/http"
	"strings"

	"github.com/denisbrodbeck/striphtmltags"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"mail/internal/models/response"
	"mail/internal/pkg/utils/validators"

	apiModels "mail/internal/models/delivery_models"
)

// AddDraft adds a new draft message.
// @Summary AddDraft a new draft message
// @Description AddDraft a new draft message to the system
// @Tags drafts-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param draft body response.EmailOtherSwag true "Draft message in JSON format"
// @Success 200 {object} response.Response "ID of the draft message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add draft message"
// @Router /api/v1/gmail/draft/adddraft [post]
func (g *GMailEmailHandler) AddDraft(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newDraft apiModels.OtherEmail
	if err := newDraft.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newDraft.Topic = sanitizeString(newDraft.Topic)
	newDraft.Text = sanitizeString(newDraft.Text)
	newDraft.SenderEmail = sanitizeString(newDraft.SenderEmail)
	newDraft.RecipientEmail = sanitizeString(newDraft.RecipientEmail)

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = g.Sessions.CheckLogin(newDraft.SenderEmail, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	strRaw := fmt.Sprintf("From: %v\r\n", newDraft.SenderEmail)
	if !validators.IsEmpty(newDraft.RecipientEmail) {
		strRaw += fmt.Sprintf("To: %v\r\n", newDraft.RecipientEmail)
	}
	strRaw += fmt.Sprintf("Subject: %v\r\n\r\n%v", newDraft.Topic, newDraft.Text)

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

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newDraft})
}

// SendDraft adds a new sent draft message.
// @Summary SendDraft a new sent draft message
// @Description SendDraft a new sent draft message to the system
// @Tags drafts-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param draft body response.EmailOtherSwag true "Draft message in JSON format"
// @Success 200 {object} response.Response "ID of the send draft message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to send draft message"
// @Router /api/v1/gmail/draft/sendDraft [post]
func (g *GMailEmailHandler) SendDraft(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newDraft apiModels.OtherEmail
	if err := newDraft.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newDraft.Topic = sanitizeString(newDraft.Topic)
	newDraft.Text = sanitizeString(newDraft.Text)
	newDraft.SenderEmail = sanitizeString(newDraft.SenderEmail)
	newDraft.RecipientEmail = sanitizeString(newDraft.RecipientEmail)

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = g.Sessions.CheckLogin(newDraft.SenderEmail, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	if validators.IsEmpty(newDraft.RecipientEmail) {
		response.HandleError(w, http.StatusBadRequest, "Empty login recipient")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	input := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("From: %v\r\nTo: %v\r\nSubject: %v\r\n\r\n%v", newDraft.SenderEmail, newDraft.RecipientEmail, newDraft.Topic, newDraft.Text)))
	draft := &gmail.Draft{
		Id: newDraft.ID,
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

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newDraft})
}

// GetByIdDraft returns an draft message by its ID.
// @Summary Get a draft message by ID
// @Description Get a draft message by its unique identifier
// @Tags drafts-gmail
// @Produce json
// @Param id path string true "ID of the draft message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Draft message data"
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

	draftResult := CreateEmailStructDraft(draft)

	for _, l := range draft.Message.LabelIds {
		label, err := srv.Users.Labels.Get("me", l).Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed get label")
			return
		}
		if label.Name == "UNREAD" {
			draftResult.ReadStatus = false
			break
		} else {
			draftResult.ReadStatus = true
		}
	}
	draftResult.DraftStatus = true

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": draftResult})
}

// UpdateDraft update a draft message.
// @Summary UpdateDraft an update draft message
// @Description UpdateDraft an update draft message to the system
// @Tags drafts-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path string true "ID of the draft message"
// @Param draft body response.EmailOtherSwag true "Draft message in JSON format"
// @Success 200 {object} response.Response "ID of the update draft message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to update draft message"
// @Router /api/v1/gmail/draft/update/{id} [put]
func (g *GMailEmailHandler) UpdateDraft(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	draftID, ok := vars["id"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newDraft apiModels.OtherEmail
	if err := newDraft.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newDraft.Topic = sanitizeString(newDraft.Topic)
	newDraft.Text = sanitizeString(newDraft.Text)
	newDraft.SenderEmail = sanitizeString(newDraft.SenderEmail)
	newDraft.RecipientEmail = sanitizeString(newDraft.RecipientEmail)

	login, err := g.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = g.Sessions.CheckLogin(newDraft.SenderEmail, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	srv, err := GetSRV(login)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	strRaw := fmt.Sprintf("From: %v\r\n", newDraft.SenderEmail)
	if !validators.IsEmpty(newDraft.RecipientEmail) {
		strRaw += fmt.Sprintf("To: %v\r\n", newDraft.RecipientEmail)
	}
	strRaw += fmt.Sprintf("Subject: %v\r\n\r\n%v", newDraft.Topic, newDraft.Text)

	input := base64.RawStdEncoding.EncodeToString([]byte(strRaw))

	draft := &gmail.Draft{
		Id: draftID,
		Message: &gmail.Message{
			Raw: input,
		},
	}

	updateDraft, err := srv.Users.Drafts.Update("me", draftID, draft).Do()
	if err != nil {
		fmt.Println(err)
		response.HandleError(w, http.StatusInternalServerError, "Error update the draft")
		return
	}

	var modifyRequest *gmail.ModifyMessageRequest
	if !newDraft.ReadStatus {
		modifyRequest = &gmail.ModifyMessageRequest{
			AddLabelIds: []string{"UNREAD"},
		}
	} else {
		modifyRequest = &gmail.ModifyMessageRequest{
			RemoveLabelIds: []string{"UNREAD"},
		}
	}
	_, err = srv.Users.Messages.Modify("me", updateDraft.Message.Id, modifyRequest).Do()
	if err != nil {
		fmt.Println(err)
		response.HandleError(w, http.StatusInternalServerError, "Error update read/unread the draft")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": newDraft})
}

// GetDrafts displays the list of draft messages.
// @Summary Display the list of draft messages
// @Description Get a list of all draft messages
// @Tags drafts-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all draft messages"
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

	p := bluemonday.StripTagsPolicy()
	draftsApi := make([]*apiModels.OtherEmail, len(req.Drafts))
	for i, d := range req.Drafts {
		dr, err := srv.Users.Drafts.Get("me", d.Id).Format("full").Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error receiving messages")
			return
		}
		email := CreateEmailStructDraft(dr)
		text := p.Sanitize(email.Text)
		text = strings.ReplaceAll(text, "\n", "")
		fields := strings.Fields(text)
		email.Text = strings.Join(fields, " ")
		for _, l := range dr.Message.LabelIds {
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
		email.DraftStatus = true
		draftsApi[i] = email
	}

	for i := range draftsApi {
		draftsApi[i].DraftStatus = true
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": draftsApi})
}

// DeleteDraft deletes an draft message.
// @Summary DeleteDraft an draft message
// @Description DeleteDraft an draft message based on its identifier
// @Tags drafts-gmail
// @Produce json
// @Param id path string true "ID of the draft message"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Deletion success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete draft message"
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
		email = ParserMessageHeadersDraft(email, msg)
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
		email = ParserMessageHeadersDraft(email, msg)
	} else if len(msg.Message.Payload.Parts) != 0 {
		for _, part := range msg.Message.Payload.Parts {
			if part.MimeType == "text/html" {
				data, err := base64.URLEncoding.DecodeString(part.Body.Data)
				if err != nil {
					fmt.Println("Error: ", err)
				}
				email.Text = string(data)
				email = ParserMessageHeadersDraft(email, msg)
			}
		}
	}

	return email
}

func ParserMessageHeadersDraft(email *apiModels.OtherEmail, msg *gmail.Draft) *apiModels.OtherEmail {
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
