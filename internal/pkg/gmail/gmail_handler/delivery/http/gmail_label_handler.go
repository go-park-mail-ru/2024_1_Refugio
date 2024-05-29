package http

import (
	"fmt"
	"google.golang.org/api/gmail/v1"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"

	"mail/internal/models/response"
	"mail/internal/pkg/utils/validators"

	apiModels "mail/internal/models/delivery_models"
)

// GetLabels displays the list of label.
// @Summary Display the list of labels
// @Description Get a list of all labels
// @Tags labels-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all labels"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/gmail/labels [get]
func (g *GMailEmailHandler) GetLabels(w http.ResponseWriter, r *http.Request) {
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

	req, err := srv.Users.Labels.List("me").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	var labelsApi []*apiModels.OtherLabel
	for _, l := range req.Labels {
		if strings.Contains(l.Id, "Label") {
			label := &apiModels.OtherLabel{
				ID:   l.Id,
				Name: l.Name,
			}
			labelsApi = append(labelsApi, label)
		}
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"labels": labelsApi})
}

// GetAllEmailsInLabel displays the list of emails in label.
// @Summary Display the list of emails in label
// @Description Get a list of all emails in label
// @Tags labels-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param name path string true "Name of the label message"
// @Success 200 {object} response.Response "List of all emails in label"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/gmail/label/{name}/emails [get]
func (g *GMailEmailHandler) GetAllEmailsInLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	labelName, ok := vars["name"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad name in request")
		return
	}

	labelName = sanitizeString(labelName)

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

	req, err := srv.Users.Messages.List("me").Q(fmt.Sprintf("label:%v", labelName)).Do()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	p := bluemonday.StripTagsPolicy()
	emailsApi := make([]*apiModels.OtherEmail, len(req.Messages))
	for i, m := range req.Messages {
		msg, _ := srv.Users.Messages.Get("me", m.Id).Format("full").Do()
		email, err := CreateEmailStruct(msg)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error decoding body data")
			return
		}
		text := p.Sanitize(email.Text)
		text = strings.ReplaceAll(text, "\n", "")
		fields := strings.Fields(text)
		email.Text = strings.Join(fields, " ")
		emailsApi[i] = email
	}
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// GetAllNameLabels displays the list of label.
// @Summary Display the list of labels
// @Description Get a list of all labels
// @Tags labels-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path string true "ID of the email message"
// @Success 200 {object} response.Response "List of all labels"
// @Failure 400 {object} response.Response "Bad request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/gmail/labels/email/{id} [get]
func (g *GMailEmailHandler) GetAllNameLabels(w http.ResponseWriter, r *http.Request) {
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

	message, err := srv.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	var labelsApi []*apiModels.OtherLabel
	for _, labelId := range message.LabelIds {
		l, err := srv.Users.Labels.Get("me", labelId).Do()
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error get label")
			return
		}
		if strings.Contains(l.Id, "Label") {
			label := &apiModels.OtherLabel{
				ID:   l.Id,
				Name: l.Name,
			}
			labelsApi = append(labelsApi, label)
		}
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"labels": labelsApi})

}

// CreateLabel adds a new label.
// @Summary CreateLabel a new label
// @Description CreateLabel a new label to the system
// @Tags labels-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param label body response.FolderSwag true "Label message in JSON format"
// @Success 200 {object} response.Response "ID of the create label"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email message"
// @Router /api/v1/gmail/label/create [post]
func (g *GMailEmailHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newLabel apiModels.OtherLabel
	if err := newLabel.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newLabel.Name = sanitizeString(newLabel.Name)

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

	color := &gmail.LabelColor{
		BackgroundColor: "#000000",
		TextColor:       "#000000",
	}

	label := &gmail.Label{
		Name:  newLabel.Name,
		Color: color,
	}

	nLabel, err := srv.Users.Labels.Create("me", label).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error create label")
		return
	}

	resultLabel := apiModels.OtherLabel{
		ID:   nLabel.Id,
		Name: nLabel.Name,
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"label": resultLabel})
}

// DeleteLabel label a user.
// @Summary DeleteLabel label a user
// @Description DeleteLabel label a user
// @Tags labels-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path string true "ID of the label"
// @Success 200 {object} response.Response "Deletion success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete label"
// @Router /api/v1/gmail/label/delete/{id} [delete]
func (g *GMailEmailHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	labelID, ok := vars["id"]
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

	err = srv.Users.Labels.Delete("me", labelID).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error delete label")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": true})
}

// UpdateLabel label a user.
// @Summary UpdateLabel label a user
// @Description UpdateLabel label a user
// @Tags labels-gmail
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path string true "ID of the label message"
// @Param label body response.FolderSwag true "Label message in JSON format"
// @Success 200 {object} response.Response "Update success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to update label"
// @Router /api/v1/gmail/label/update/{id} [put]
func (g *GMailEmailHandler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	labelID, ok := vars["id"]
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newLabel apiModels.OtherLabel
	if err := newLabel.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
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

	color := &gmail.LabelColor{
		BackgroundColor: "#000000",
		TextColor:       "#000000",
	}

	Label := &gmail.Label{
		Name:  newLabel.Name,
		Color: color,
	}

	_, err = srv.Users.Labels.Update("me", labelID, Label).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error update label")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": true})
}

// AddEmailInLabel adds an email in label.
// @Summary AddEmailInLabel an email in label
// @Description AddEmailInLabel an email in label to the system
// @Tags labels-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param label body response.FolderEmailGoogleSwag true "Label message in JSON format"
// @Success 200 {object} response.Response "Success of the add email in label"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email in label"
// @Router /api/v1/gmail/label/add_email [post]
func (g *GMailEmailHandler) AddEmailInLabel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newLabelEmail apiModels.LabelEmail
	if err := newLabelEmail.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
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

	modifyRequest := &gmail.ModifyMessageRequest{
		AddLabelIds: []string{newLabelEmail.LabelID},
	}

	_, err = srv.Users.Messages.Modify("me", newLabelEmail.EmailID, modifyRequest).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to add LabelEmail message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": true})
}

// DeleteEmailInLabel delete email in label.
// @Summary DeleteEmailInLabel an email in label
// @Description DeleteEmailInLabel an email in label to the system
// @Tags labels-gmail
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param label body response.FolderEmailGoogleSwag true "Label message in JSON format"
// @Success 200 {object} response.Response "Success of the delete email in label"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete email in label"
// @Router /api/v1/gmail/label/delete_email [delete]
func (g *GMailEmailHandler) DeleteEmailInLabel(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newLabelEmail apiModels.LabelEmail
	if err := newLabelEmail.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
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

	modifyRequest := &gmail.ModifyMessageRequest{
		RemoveLabelIds: []string{newLabelEmail.LabelID},
	}

	_, err = srv.Users.Messages.Modify("me", newLabelEmail.EmailID, modifyRequest).Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete LabelEmail message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": true})
}
