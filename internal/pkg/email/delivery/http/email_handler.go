package http

import (
	"encoding/json"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	emailUsecase "mail/internal/microservice/email/interface"
	"mail/internal/microservice/email/proto"
	"mail/internal/microservice/models/proto_converters"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"
	"net/http"
	"strconv"

	converters "mail/internal/models/delivery_converters"

	emailApi "mail/internal/models/delivery_models"
	"mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"mail/internal/pkg/utils/validators"

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

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	emailServiceClient := proto.NewEmailServiceClient(conn)
	emailDataProto, err := emailServiceClient.GetAllIncoming(
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

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	emailServiceClient := proto.NewEmailServiceClient(conn)
	emailDataProto, err := emailServiceClient.GetAllSent(
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

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	emailServiceClient := proto.NewEmailServiceClient(conn)
	emailDataProto, err := emailServiceClient.GetEmailByID(
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
	newEmail.AvatarID = sanitizeString(newEmail.AvatarID)
	newEmail.SenderEmail = sanitizeString(newEmail.SenderEmail)
	newEmail.RecipientEmail = sanitizeString(newEmail.RecipientEmail)

	sender := newEmail.SenderEmail
	recipient := newEmail.RecipientEmail

	switch {
	case validators.IsValidEmailFormat(sender) && validators.IsValidEmailFormat(recipient):
		err = h.Sessions.CheckLogin(sender, r, r.Context())
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad sender login")
			return
		}

		conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer conn.Close()

		emailServiceClient := proto.NewEmailServiceClient(conn)
		_, err = emailServiceClient.CheckRecipientEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Recipient{Recipient: recipient},
		)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad login")
			return
		}

		emailDataProto, err := emailServiceClient.CreateEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Email{
				Id:             newEmail.ID,
				Topic:          newEmail.Topic,
				Text:           newEmail.Text,
				PhotoID:        newEmail.AvatarID,
				ReadStatus:     newEmail.ReadStatus,
				Flag:           newEmail.Flag,
				Deleted:        newEmail.Deleted,
				DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
				ReplyToEmailID: newEmail.ReplyToEmailID,
				DraftStatus:    newEmail.DraftStatus,
				SpamStatus:     newEmail.SpamStatus,
				SenderEmail:    newEmail.SenderEmail,
				RecipientEmail: newEmail.RecipientEmail,
			},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}
		emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
		emailData.ID = emailDataProto.Id

		_, err = emailServiceClient.CreateProfileEmail(
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
		/*email_id, email, err := h.EmailUseCase.CreateEmail(converters.EmailConvertApiInCore(newEmail))
		if err != nil {
			delivery.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}

		delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"email": converters.EmailConvertCoreInApi(*email)})*/
		response.HandleSuccess(w, http.StatusBadRequest, "An error occurred in the recipient's domain. You cannot send messages to other email services. Make sure that the recipient's domain ends with @mailhub.su")
		return
	case validators.IsValidEmailFormat(sender) == false && validators.IsValidEmailFormat(recipient) == true:
		conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer conn.Close()

		emailServiceClient := proto.NewEmailServiceClient(conn)
		_, err = emailServiceClient.CheckRecipientEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Recipient{Recipient: recipient},
		)
		if err != nil {
			response.HandleError(w, http.StatusBadRequest, "Bad login")
			return
		}

		emailDataProto, err := emailServiceClient.CreateEmail(
			metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
			&proto.Email{
				Id:             newEmail.ID,
				Topic:          newEmail.Topic,
				Text:           newEmail.Text,
				PhotoID:        newEmail.AvatarID,
				ReadStatus:     newEmail.ReadStatus,
				Flag:           newEmail.Flag,
				Deleted:        newEmail.Deleted,
				DateOfDispatch: timestamppb.New(newEmail.DateOfDispatch),
				ReplyToEmailID: newEmail.ReplyToEmailID,
				DraftStatus:    newEmail.DraftStatus,
				SpamStatus:     newEmail.SpamStatus,
				SenderEmail:    newEmail.SenderEmail,
				RecipientEmail: newEmail.RecipientEmail,
			},
		)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "Failed to add email message")
			return
		}
		emailData := proto_converters.EmailConvertProtoInCore(*emailDataProto.Email)
		emailData.ID = emailDataProto.Id

		_, err = emailServiceClient.CreateProfileEmail(
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
	updatedEmail.AvatarID = sanitizeString(updatedEmail.AvatarID)
	updatedEmail.RecipientEmail = sanitizeString(updatedEmail.RecipientEmail)
	updatedEmail.SenderEmail = sanitizeString(updatedEmail.SenderEmail)

	updatedEmail.ID = id

	err1 := h.Sessions.CheckLogin(updatedEmail.SenderEmail, r, r.Context())
	err2 := h.Sessions.CheckLogin(updatedEmail.RecipientEmail, r, r.Context())
	if err1 != nil && err2 != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad sender login")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	emailServiceClient := proto.NewEmailServiceClient(conn)
	emailDataProto, err := emailServiceClient.UpdateEmail(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.Email{
			Id:             updatedEmail.ID,
			Topic:          updatedEmail.Topic,
			Text:           updatedEmail.Text,
			PhotoID:        updatedEmail.AvatarID,
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

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	emailServiceClient := proto.NewEmailServiceClient(conn)
	emailDataProto, err := emailServiceClient.DeleteEmail(
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

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	emailServiceClient := proto.NewEmailServiceClient(conn)
	emailDataProto, err := emailServiceClient.GetDraftEmails(
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

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.EmailService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	emailServiceClient := proto.NewEmailServiceClient(conn)
	emailDataProto, err := emailServiceClient.GetSpamEmails(
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
