package http

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/minio/minio-go/v7"

	"mail/cmd/configs"
	"mail/internal/microservice/email/proto"
	"mail/internal/microservice/models/proto_converters"
	"mail/internal/models/response"
	"mail/internal/pkg/utils/check_file_type"
	"mail/internal/pkg/utils/generate_filename"
	"mail/internal/pkg/utils/validators"

	email_proto "mail/internal/microservice/email/proto"
	converters "mail/internal/models/delivery_converters"
	emailApi "mail/internal/models/delivery_models"
	domainSession "mail/internal/pkg/session/interface"
)

var (
	EHandler                        = &EmailHandler{}
	requestIDContextKey interface{} = "requestid"
)

// EmailHandler represents the handler for email operations.
type EmailHandler struct {
	Sessions           domainSession.SessionsManager
	EmailServiceClient email_proto.EmailServiceClient
	MinioClient        *minio.Client
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
	login, err := h.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user session")
		return
	}

	err = h.Sessions.CheckLogin(login, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user login")
		return
	}

	emailDataProto, err := h.EmailServiceClient.GetAllIncoming(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.LoginOffsetLimit{Login: login, Offset: 0, Limit: 0},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	emailsCore := proto_converters.EmailsConvertProtoInCore(emailDataProto)

	emailsApi := make([]*emailApi.Email, 0, len(emailsCore))
	for _, email := range emailsCore {
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
	login, err := h.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user session")
		return
	}

	err = h.Sessions.CheckLogin(login, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user login")
		return
	}

	emailDataProto, err := h.EmailServiceClient.GetAllSent(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.LoginOffsetLimit{Login: login, Offset: 0, Limit: 0},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	emailsCore := proto_converters.EmailsConvertProtoInCore(emailDataProto)

	emailsApi := make([]*emailApi.Email, 0, len(emailsCore))
	for _, email := range emailsCore {
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

	login, err := h.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = h.Sessions.CheckLogin(login, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	emailDataProto, err := h.EmailServiceClient.GetEmailByID(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.EmailIdAndLogin{Id: id, Login: login},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, "Email not found")
		return
	}
	emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newEmail emailApi.Email
	if err := newEmail.UnmarshalJSON(body); err != nil {
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

	if validators.IsEmpty(newEmail.Text) || validators.IsEmpty(sender) || validators.IsEmpty(recipient) {
		response.HandleError(w, http.StatusBadRequest, "Data is empty")
		return
	}

	switch {
	case validators.IsValidEmailFormat(sender) && validators.IsValidEmailFormat(recipient):
		err = h.Sessions.CheckLogin(sender, r, r.Context())
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad sender login")
			return
		}

		_, err = h.EmailServiceClient.CheckRecipientEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Recipient{Recipient: recipient},
		)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad login")
			return
		}

		emailDataProto, err := h.EmailServiceClient.CreateEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Email{
				Id:             newEmail.ID,
				Topic:          newEmail.Topic,
				Text:           newEmail.Text,
				PhotoID:        newEmail.PhotoID,
				ReadStatus:     newEmail.ReadStatus,
				Flag:           newEmail.Flag,
				Deleted:        newEmail.Deleted,
				DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
				ReplyToEmailID: newEmail.ReplyToEmailID,
				DraftStatus:    newEmail.DraftStatus,
				SpamStatus:     newEmail.SpamStatus,
				SenderEmail:    sender,
				RecipientEmail: recipient,
			},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}
		emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
		emailData.ID = emailDataProto.Id

		_, err = h.EmailServiceClient.CreateProfileEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.IdSenderRecipient{Id: emailData.ID, Sender: emailData.SenderEmail, Recipient: emailData.RecipientEmail},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
		return
	case validators.IsValidEmailFormat(sender) == true && validators.IsValidEmailFormat(recipient) == false:
		err = h.Sessions.CheckLogin(sender, r, r.Context())
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad sender login")
			return
		}

		var mx string
		var err error

		// Connect to the server, authenticate, set the sender and recipient,
		// and send the email all in one step.
		//from := "sergey@mailhub.su"
		//to := "fedasovsergey00@gmail.com"
		//subject := "Test Email"
		//body := "This is the email body at " + time.Now().String()

		msg := composeMimeMail(recipient, sender, newEmail.Topic, newEmail.Text)

		mx, err = getMXRecord(recipient)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad request to login")
			return
		}

		err = smtp.SendMail(mx+":25", nil, sender, []string{recipient}, msg)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Error send email")
			return
		}

		emailDataProto, err := h.EmailServiceClient.CreateEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Email{
				Id:             newEmail.ID,
				Topic:          newEmail.Topic,
				Text:           newEmail.Text,
				PhotoID:        newEmail.PhotoID,
				ReadStatus:     newEmail.ReadStatus,
				Flag:           newEmail.Flag,
				Deleted:        newEmail.Deleted,
				DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
				ReplyToEmailID: newEmail.ReplyToEmailID,
				DraftStatus:    newEmail.DraftStatus,
				SpamStatus:     newEmail.SpamStatus,
				SenderEmail:    sender,
				RecipientEmail: recipient,
			},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}
		emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
		emailData.ID = emailDataProto.Id

		_, err = h.EmailServiceClient.CreateProfileEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.IdSenderRecipient{Id: emailData.ID, Sender: emailData.SenderEmail, Recipient: emailData.RecipientEmail},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
		return
	case validators.IsValidEmailFormat(sender) == false && validators.IsValidEmailFormat(recipient) == true:
		fmt.Println("-----START_1----")
		_, err = h.EmailServiceClient.CheckRecipientEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Recipient{Recipient: recipient},
		)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad login")
			return
		}

		fmt.Println("-----START_2----")
		fmt.Println(newEmail)
		emailDataProto, err := h.EmailServiceClient.CreateEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Email{
				Id:             newEmail.ID,
				Topic:          newEmail.Topic,
				Text:           newEmail.Text,
				PhotoID:        newEmail.PhotoID,
				ReadStatus:     newEmail.ReadStatus,
				Flag:           newEmail.Flag,
				Deleted:        newEmail.Deleted,
				DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
				ReplyToEmailID: newEmail.ReplyToEmailID,
				DraftStatus:    newEmail.DraftStatus,
				SpamStatus:     newEmail.SpamStatus,
				SenderEmail:    sender,
				RecipientEmail: recipient,
			},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}
		emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
		emailData.ID = emailDataProto.Id

		fmt.Println("-----START_3----")
		_, err = h.EmailServiceClient.CreateProfileEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.IdSenderRecipient{Id: emailData.ID, Sender: emailData.SenderEmail, Recipient: emailData.RecipientEmail},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
		return
	}
}

func getMXRecord(to string) (mx string, err error) {
	var e *mail.Address
	e, err = mail.ParseAddress(to)
	if err != nil {
		return
	}

	domain := strings.Split(e.Address, "@")[1]

	var mxs []*net.MX
	mxs, err = net.LookupMX(domain)

	if err != nil {
		return
	}

	for _, x := range mxs {
		mx = x.Host
		return
	}

	return
}

func formatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Address: str}
	return strings.Trim(addr.String(), " <>")
}

func composeMimeMail(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = formatEmailAddress(from)
	header["To"] = formatEmailAddress(to)
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var updatedEmail emailApi.Email
	if err := updatedEmail.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	updatedEmail.Topic = sanitizeString(updatedEmail.Topic)
	updatedEmail.Text = sanitizeString(updatedEmail.Text)
	updatedEmail.PhotoID = sanitizeString(updatedEmail.PhotoID)
	updatedEmail.RecipientEmail = sanitizeString(updatedEmail.RecipientEmail)
	updatedEmail.SenderEmail = sanitizeString(updatedEmail.SenderEmail)

	updatedEmail.ID = id

	err1 := h.Sessions.CheckLogin(updatedEmail.SenderEmail, r, r.Context())
	err2 := h.Sessions.CheckLogin(updatedEmail.RecipientEmail, r, r.Context())
	if err1 != nil && err2 != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	emailDataProto, err := h.EmailServiceClient.UpdateEmail(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.Email{
			Id:             updatedEmail.ID,
			Topic:          updatedEmail.Topic,
			Text:           updatedEmail.Text,
			PhotoID:        updatedEmail.PhotoID,
			ReadStatus:     updatedEmail.ReadStatus,
			Flag:           updatedEmail.Flag,
			Deleted:        updatedEmail.Deleted,
			DateOfDispatch: timestamppb.New(updatedEmail.DateOfDispatch),
			ReplyToEmailID: updatedEmail.ReplyToEmailID,
			DraftStatus:    updatedEmail.DraftStatus,
			SpamStatus:     updatedEmail.SpamStatus,
			SenderEmail:    updatedEmail.SenderEmail,
			RecipientEmail: updatedEmail.RecipientEmail,
		},
	)
	if err != nil || !emailDataProto.Status {
		response.HandleError(w, http.StatusInternalServerError, "Failed to update email message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": emailDataProto.Status})
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

	login, err := h.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender session")
		return
	}

	err = h.Sessions.CheckLogin(login, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	emailDataProto, err := h.EmailServiceClient.DeleteEmail(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.LoginWithID{Id: id, Login: login},
	)
	if err != nil || !emailDataProto.Status {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete email message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": emailDataProto.Status})
}

// Draft displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 404 {object} response.Response "DB error"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/emails/draft [get]
func (h *EmailHandler) Draft(w http.ResponseWriter, r *http.Request) {
	login, err := h.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user session")
		return
	}

	err = h.Sessions.CheckLogin(login, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user login")
		return
	}

	emailDataProto, err := h.EmailServiceClient.GetDraftEmails(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.LoginOffsetLimit{Login: login, Offset: 0, Limit: 0},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	emailsCore := proto_converters.EmailsConvertProtoInCore(emailDataProto)

	emailsApi := make([]*emailApi.Email, 0, len(emailsCore))
	for _, email := range emailsCore {
		emailsApi = append(emailsApi, converters.EmailConvertCoreInApi(*email))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

// Spam displays the list of email messages.
// @Summary Display the list of email messages
// @Description Get a list of all email messages
// @Tags emails
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "List of all email messages"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 404 {object} response.Response "DB error"
// @Failure 500 {object} response.Response "JSON encoding error"
// @Router /api/v1/emails/spam [get]
func (h *EmailHandler) Spam(w http.ResponseWriter, r *http.Request) {
	login, err := h.Sessions.GetLoginBySession(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user session")
		return
	}

	err = h.Sessions.CheckLogin(login, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad user login")
		return
	}

	emailDataProto, err := h.EmailServiceClient.GetSpamEmails(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.LoginOffsetLimit{Login: login, Offset: 0, Limit: 0},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("DB error: %s", err.Error()))
		return
	}

	emailsCore := proto_converters.EmailsConvertProtoInCore(emailDataProto)

	emailsApi := make([]*emailApi.Email, 0, len(emailsCore))
	for _, email := range emailsCore {
		emailsApi = append(emailsApi, converters.EmailConvertCoreInApi(*email))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"emails": emailsApi})
}

func (h *EmailHandler) SendFromAnotherDomain(w http.ResponseWriter, r *http.Request) {
	h.Send(w, r)
}

// AddDraft adds a new draft email message.
// @Summary AddDraft a new draft email message
// @Description AddDraft a new draft email message to the system
// @Tags emails
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param email body response.EmailSwag true "Email message in JSON format"
// @Success 200 {object} response.Response "ID of the send email message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add email message"
// @Router /api/v1/email/adddraft [post]
func (h *EmailHandler) AddDraft(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newEmail emailApi.Email
	if err := newEmail.UnmarshalJSON(body); err != nil {
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

	if validators.IsEmpty(sender) {
		response.HandleError(w, http.StatusBadRequest, "Data is empty")
		return
	}

	if validators.IsEmpty(recipient) && validators.IsValidEmailFormat(sender) {
		err = h.Sessions.CheckLogin(sender, r, r.Context())
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad sender login")
			return
		}

		emailDataProto, err := h.EmailServiceClient.AddEmailDraft(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Email{
				Id:             newEmail.ID,
				Topic:          newEmail.Topic,
				Text:           newEmail.Text,
				PhotoID:        newEmail.PhotoID,
				ReadStatus:     newEmail.ReadStatus,
				Flag:           newEmail.Flag,
				Deleted:        newEmail.Deleted,
				DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
				ReplyToEmailID: newEmail.ReplyToEmailID,
				DraftStatus:    newEmail.DraftStatus,
				SpamStatus:     newEmail.SpamStatus,
				SenderEmail:    sender,
				RecipientEmail: recipient,
			},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}
		emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
		emailData.ID = emailDataProto.Id

		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
		return
	} else {
		switch {
		case validators.IsValidEmailFormat(sender) && validators.IsValidEmailFormat(recipient):
			err = h.Sessions.CheckLogin(sender, r, r.Context())
			if err != nil {
				response.HandleError(w, http.StatusBadRequest, "Bad sender login")
				return
			}

			_, err = h.EmailServiceClient.CheckRecipientEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.Recipient{Recipient: recipient},
			)
			if err != nil {
				response.HandleError(w, http.StatusBadRequest, "Bad login")
				return
			}

			emailDataProto, err := h.EmailServiceClient.CreateEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.Email{
					Id:             newEmail.ID,
					Topic:          newEmail.Topic,
					Text:           newEmail.Text,
					PhotoID:        newEmail.PhotoID,
					ReadStatus:     newEmail.ReadStatus,
					Flag:           newEmail.Flag,
					Deleted:        newEmail.Deleted,
					DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
					ReplyToEmailID: newEmail.ReplyToEmailID,
					DraftStatus:    newEmail.DraftStatus,
					SpamStatus:     newEmail.SpamStatus,
					SenderEmail:    sender,
					RecipientEmail: recipient,
				},
			)
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
				return
			}
			emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
			emailData.ID = emailDataProto.Id

			_, err = h.EmailServiceClient.CreateProfileEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.IdSenderRecipient{Id: emailData.ID, Sender: emailData.SenderEmail, Recipient: emailData.RecipientEmail},
			)
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
				return
			}

			response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
			return
		case validators.IsValidEmailFormat(sender) == true && validators.IsValidEmailFormat(recipient) == false:
			err = h.Sessions.CheckLogin(sender, r, r.Context())
			if err != nil {
				response.HandleError(w, http.StatusBadRequest, "Bad sender login")
				return
			}

			emailDataProto, err := h.EmailServiceClient.CreateEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.Email{
					Id:             newEmail.ID,
					Topic:          newEmail.Topic,
					Text:           newEmail.Text,
					PhotoID:        newEmail.PhotoID,
					ReadStatus:     newEmail.ReadStatus,
					Flag:           newEmail.Flag,
					Deleted:        newEmail.Deleted,
					DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
					ReplyToEmailID: newEmail.ReplyToEmailID,
					DraftStatus:    newEmail.DraftStatus,
					SpamStatus:     newEmail.SpamStatus,
					SenderEmail:    sender,
					RecipientEmail: recipient,
				},
			)
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
				return
			}
			emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
			emailData.ID = emailDataProto.Id

			_, err = h.EmailServiceClient.CreateProfileEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.IdSenderRecipient{Id: emailData.ID, Sender: emailData.SenderEmail, Recipient: emailData.RecipientEmail},
			)
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
				return
			}

			response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
			return
		case validators.IsValidEmailFormat(sender) == false && validators.IsValidEmailFormat(recipient) == true:
			_, err = h.EmailServiceClient.CheckRecipientEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.Recipient{Recipient: recipient},
			)
			if err != nil {
				response.HandleError(w, http.StatusBadRequest, "Bad login")
				return
			}

			emailDataProto, err := h.EmailServiceClient.CreateEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.Email{
					Id:             newEmail.ID,
					Topic:          newEmail.Topic,
					Text:           newEmail.Text,
					PhotoID:        newEmail.PhotoID,
					ReadStatus:     newEmail.ReadStatus,
					Flag:           newEmail.Flag,
					Deleted:        newEmail.Deleted,
					DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
					ReplyToEmailID: newEmail.ReplyToEmailID,
					DraftStatus:    newEmail.DraftStatus,
					SpamStatus:     newEmail.SpamStatus,
					SenderEmail:    sender,
					RecipientEmail: recipient,
				},
			)
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
				return
			}
			emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
			emailData.ID = emailDataProto.Id

			_, err = h.EmailServiceClient.CreateProfileEmail(
				metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
				&proto.IdSenderRecipient{Id: emailData.ID, Sender: emailData.SenderEmail, Recipient: emailData.RecipientEmail},
			)
			if err != nil {
				response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
				return
			}

			response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*emailData)})
			return
		}
	}
}

// AddAttachment adds an attachment to an email message.
// @Summary Add an attachment to an email message
// @Description Add an attachment to an email message by its ID
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param id path integer true "ID of the email message"
// @Param file formData file true "Attachment file to upload"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Attachment added successfully"
// @Failure 400 {object} response.Response "Bad id in request or bad JSON in request"
// @Failure 404 {object} response.Response "Failed to add attachment"
// @Router /api/v1/email/{id}/addattachment [post]
func (h *EmailHandler) AddAttachment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	err = r.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Error processing file")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	if handler.Size > (20 * 1024 * 1024) {
		response.HandleError(w, http.StatusInternalServerError, "Failed to get file")
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	contentType := handler.Header.Get("Content-Type")

	fileType := check_file_type.GetFileType(contentType)

	uniqueFileName := generate_filename.GenerateUniqueFileName(fileExt)

	_, err = h.MinioClient.PutObject(r.Context(), "files", uniqueFileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error uploading file to MinIO")
		return
	}

	fileURL := fmt.Sprintf(configs.PROTOCOL+"0.0.0.0:9000"+"/files/%s", uniqueFileName)

	fileId, err := h.EmailServiceClient.AddAttachment(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.AddAttachmentRequest{EmailId: id, FileId: fileURL, FileType: fileType},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to add attachment: %s", err.Error()))
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"FileId": fileId.FileId})
}

// GetFileByID retrieves a file by its ID.
// @Summary Retrieve a file by its ID
// @Description Retrieve a file by its ID from the email service
// @Tags files
// @Produce json
// @Param id path integer true "ID of the file"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "File retrieved successfully"
// @Failure 400 {object} response.Response "Bad ID in request"
// @Failure 404 {object} response.Response "File not found"
// @Router /api/v1/email/get/file/{id} [get]
func (h *EmailHandler) GetFileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad ID in request")
		return
	}

	fileProto, err := h.EmailServiceClient.GetFileByID(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.GetFileByIDRequest{FileId: id},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to get file: %s", err.Error()))
		return
	}

	fileApi := emailApi.File{
		ID:       fileProto.File.Id,
		FileId:   fileProto.File.FileId,
		FileType: fileProto.File.FileType,
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"file": fileApi})
}

// GetFilesByEmailID retrieves files associated with an email by its ID.
// @Summary Retrieve files by email ID
// @Description Retrieve files associated with an email by its ID from the email service
// @Tags files
// @Produce json
// @Param id path integer true "ID of the email"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Files retrieved successfully"
// @Failure 400 {object} response.Response "Bad ID in request"
// @Failure 404 {object} response.Response "Files not found"
// @Router /api/v1/email/{id}/get/files/ [get]
func (h *EmailHandler) GetFilesByEmailID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad ID in request")
		return
	}

	filesProto, err := h.EmailServiceClient.GetFilesByEmailID(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.GetFilesByEmailIDRequest{EmailId: id},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to get files: %s", err.Error()))
		return
	}

	filesApi := make([]*emailApi.File, 0, len(filesProto.Files))
	for _, file := range filesProto.Files {
		filesApi = append(filesApi, &emailApi.File{
			ID:       file.Id,
			FileId:   file.FileId,
			FileType: file.FileType,
		})
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"files": filesApi})
}

// DeleteFileByID deletes a file by its ID.
// @Summary Delete a file by its ID
// @Description Delete a file by its ID from the email service
// @Tags files
// @Produce json
// @Param id path integer true "ID of the file"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "File deleted successfully"
// @Failure 400 {object} response.Response "Bad ID in request"
// @Failure 404 {object} response.Response "Failed to delete file"
// @Router /api/v1/email/delete/file/{id} [delete]
func (h *EmailHandler) DeleteFileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad ID in request")
		return
	}

	status, err := h.EmailServiceClient.DeleteFileByID(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.DeleteFileByIDRequest{FileId: id},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to delete file: %s", err.Error()))
		return
	}
	if !status.Status {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to delete file: %s", err.Error()))
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": status.Status})
}

// UpdateFileByID updates a file by its ID.
// @Summary Update a file by its ID
// @Description Update a file by its ID in the email service
// @Tags files
// @Produce json
// @Param id path integer true "ID of the file"
// @Param file formData file true "Updated file to upload"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "File updated successfully"
// @Failure 400 {object} response.Response "Bad ID in request"
// @Failure 404 {object} response.Response "Failed to update file"
// @Router /api/v1/email/update/file/{id} [put]
func (h *EmailHandler) UpdateFileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad ID in request")
		return
	}

	err = r.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Error processing file")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	if handler.Size > (20 * 1024 * 1024) {
		response.HandleError(w, http.StatusInternalServerError, "Failed to get file")
		return
	}

	fileProto, err := h.EmailServiceClient.GetFileByID(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.GetFileByIDRequest{FileId: id},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to get file: %s", err.Error()))
		return
	}

	err = h.MinioClient.RemoveObject(r.Context(), "files", fileProto.File.FileId[strings.LastIndex(fileProto.File.FileId, "/")+1:], minio.RemoveObjectOptions{})
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete old file in MinIO")
		return
	}

	_, err = h.MinioClient.PutObject(r.Context(), "files", fileProto.File.FileId[strings.LastIndex(fileProto.File.FileId, "/")+1:], file, -1, minio.PutObjectOptions{})
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to update file in MinIO")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": true})
}

// AddFile add a file to an email message.
// @Summary Add a file to an email message
// @Description Add a file to an email message
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Attachment file to upload"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "File added successfully"
// @Failure 400 {object} response.Response "Bad id in request or bad JSON in request"
// @Failure 404 {object} response.Response "Failed to add file"
// @Router /api/v1/email/addfile [post]
func (h *EmailHandler) AddFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Error processing file")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	if handler.Size > (20 * 1024 * 1024) {
		response.HandleError(w, http.StatusInternalServerError, "Failed to get file")
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	contentType := handler.Header.Get("Content-Type")

	fileType := check_file_type.GetFileType(contentType)

	uniqueFileName := generate_filename.GenerateUniqueFileName(fileExt)

	_, err = h.MinioClient.PutObject(r.Context(), "files", uniqueFileName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error uploading file to MinIO")
		return
	}

	fileURL := fmt.Sprintf(configs.PROTOCOL+"0.0.0.0:9000"+"/files/%s", uniqueFileName)

	fileId, err := h.EmailServiceClient.AddFile(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.AddFileRequest{FileId: fileURL, FileType: fileType},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to add file: %s", err.Error()))
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"FileId": fileId.FileId})
}

// AddFileToEmail adds a file to an email message.
// @Summary Add a file to an email message
// @Description Adds a file as an attachment to a specified email message
// @Tags files
// @Produce json
// @Param id path uint64 true "Email ID"
// @Param file-id path uint64 true "File ID"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "File added successfully"
// @Failure 400 {object} response.Response "Bad ID in request or error processing file"
// @Failure 404 {object} response.Response "Failed to add file"
// @Router /api/v1/email/{id}/file/{file-id} [post]
func (h *EmailHandler) AddFileToEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	email_id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad ID in request")
		return
	}

	file_id, err := strconv.ParseUint(vars["file-id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad ID in request")
		return
	}

	status, err := h.EmailServiceClient.AddFileToEmail(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.AddFileToEmailRequest{EmailId: email_id, FileId: file_id},
	)
	if err != nil {
		response.HandleError(w, http.StatusNotFound, fmt.Sprintf("Failed to add file: %s", err.Error()))
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": status.Status})
}
