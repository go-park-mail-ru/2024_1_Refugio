package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc/metadata"
	folderUsecase "mail/internal/microservice/folder/interface"
	"mail/internal/microservice/folder/proto"
	"mail/internal/microservice/models/proto_converters"
	converters "mail/internal/models/delivery_converters"
	folderApi "mail/internal/models/delivery_models"
	"mail/internal/models/microservice_ports"
	"mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"mail/internal/pkg/utils/connect_microservice"
	"net/http"
	"strconv"
)

var (
	FHandler                        = &FolderHandler{}
	requestIDContextKey interface{} = "requestid"
)

// FolderHandler represents the handler for folder operations.
type FolderHandler struct {
	FolderUseCase folderUsecase.FolderUseCase
	Sessions      domainSession.SessionsManager
}

func sanitizeString(str string) string {
	p := bluemonday.UGCPolicy()
	p.AllowElements("b", "i", "a", "strong", "em", "p", "br", "span", "ul", "ol", "li", "h1", "h2", "h3", "div")
	return p.Sanitize(str)
}

// Add adds a new folder message.
// @Summary Add a new folder message
// @Description Add a new folder message to the system
// @Tags folders
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param folder body response.FolderSwag true "Folder message in JSON format"
// @Success 200 {object} response.Response "ID of the send folder message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add folder message"
// @Router /api/v1/folder/add [post]
func (h *FolderHandler) Add(w http.ResponseWriter, r *http.Request) {
	var newFolder folderApi.Folder
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newFolder)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	profileId, err := h.Sessions.GetProfileIDBySessionID(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad session")
		return
	}

	newFolder.Name = sanitizeString(newFolder.Name)
	newFolder.ProfileId = profileId

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	folderServiceClient := proto.NewFolderServiceClient(conn)
	folderDataProto, err := folderServiceClient.CreateFolder(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.Folder{
			Id:        newFolder.ID,
			ProfileId: newFolder.ProfileId,
			Name:      newFolder.Name,
		},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to add folder message")
		return
	}
	folderData := proto_converters.FolderConvertProtoInCore(*folderDataProto.Folder)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"folder": converters.FolderConvertCoreInApi(*folderData)})
}

// GetAll get all folders.
// @Summary GetAll get all folders
// @Description GetAll folders users
// @Tags folders
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "ID of the send folder message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to get all folders"
// @Router /api/v1/folder/all [get]
func (h *FolderHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	profileId, err := h.Sessions.GetProfileIDBySessionID(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad session")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	folderServiceClient := proto.NewFolderServiceClient(conn)
	folderDataProto, err := folderServiceClient.GetAllFolders(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.GetAllFoldersData{Id: profileId, Offset: 0, Limit: 0},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to get all folders")
		return
	}

	foldersCore := proto_converters.FoldersConvertProtoInCore(folderDataProto)

	foldersApi := make([]*folderApi.Folder, 0, len(foldersCore))
	for _, folder := range foldersCore {
		foldersApi = append(foldersApi, converters.FolderConvertCoreInApi(*folder))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"folders": foldersApi})
}

// Delete folder a user.
// @Summary Delete folder a user
// @Description Delete folder a user
// @Tags folders
// @Produce json
// @Param id path integer true "ID of the folder"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Deletion success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete folder"
// @Router /api/v1/folder/delete/{id} [delete]
func (h *FolderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	profileId, err := h.Sessions.GetProfileIDBySessionID(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad session")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	folderServiceClient := proto.NewFolderServiceClient(conn)
	folderDataProto, err := folderServiceClient.DeleteFolder(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.DeleteFolderData{FolderID: uint32(id), ProfileID: profileId},
	)
	if err != nil || !folderDataProto.Status {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete folder")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": folderDataProto.Status})
}

// Update folder a user.
// @Summary Update folder a user
// @Description Update folder a user
// @Tags folders
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path integer true "ID of the folder message"
// @Param folder body response.FolderSwag true "Folder message in JSON format"
// @Success 200 {object} response.Response "Update success status"
// @Failure 400 {object} response.Response "Bad id"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to update folder"
// @Router /api/v1/folder/update/{id} [put]
func (h *FolderHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	var newFolder folderApi.Folder
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err = json.NewDecoder(r.Body).Decode(&newFolder)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	profileId, err := h.Sessions.GetProfileIDBySessionID(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad session")
		return
	}

	newFolder.Name = sanitizeString(newFolder.Name)
	newFolder.ID = uint32(id)
	newFolder.ProfileId = profileId

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	folderServiceClient := proto.NewFolderServiceClient(conn)
	folderDataProto, err := folderServiceClient.UpdateFolder(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.Folder{
			Id:        newFolder.ID,
			ProfileId: newFolder.ProfileId,
			Name:      newFolder.Name,
		},
	)
	if err != nil || !folderDataProto.Status {
		response.HandleError(w, http.StatusInternalServerError, "Failed to update folder message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": folderDataProto.Status})
}

// AddEmailInFolder adds a new folder message.
// @Summary AddEmailInFolder a new folder message
// @Description AddEmailInFolder a new folder message to the system
// @Tags folders
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param folder body response.FolderEmailSwag true "Folder message in JSON format"
// @Success 200 {object} response.Response "ID of the send folder message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to add folder message"
// @Router /api/v1/folder/add_email [post]
func (h *FolderHandler) AddEmailInFolder(w http.ResponseWriter, r *http.Request) {
	var newFolderEmail folderApi.FolderEmail
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newFolderEmail)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	profileId, err := h.Sessions.GetProfileIDBySessionID(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad session")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	folderServiceClient := proto.NewFolderServiceClient(conn)
	CheckFolderProfileDataProto, err := folderServiceClient.CheckFolderProfile(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.FolderProfile{
			FolderID:  newFolderEmail.FolderID,
			ProfileID: profileId,
		},
	)
	if err != nil || !CheckFolderProfileDataProto.Status {
		response.HandleError(w, http.StatusBadRequest, "ProfileID and FolderID not found")
		return
	}

	CheckEmailProfileDataProto, err := folderServiceClient.CheckEmailProfile(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.EmailProfile{
			EmailID:   newFolderEmail.EmailID,
			ProfileID: profileId,
		},
	)
	if err != nil || !CheckEmailProfileDataProto.Status {
		response.HandleError(w, http.StatusBadRequest, "ProfileID and EmailID not found")
		return
	}

	folderDataProto, err := folderServiceClient.AddEmailInFolder(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.FolderEmail{
			FolderID: newFolderEmail.FolderID,
			EmailID:  newFolderEmail.EmailID,
		},
	)
	if err != nil || !folderDataProto.Status {
		response.HandleError(w, http.StatusInternalServerError, "Failed to add FolderEmail message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": folderDataProto.Status})
}

// DeleteEmailInFolder adds a new folder message.
// @Summary DeleteEmailInFolder a new folder message
// @Description DeleteEmailInFolder a new folder message to the system
// @Tags folders
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param folder body response.FolderEmailSwag true "Folder message in JSON format"
// @Success 200 {object} response.Response "ID of the folder message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to delete folder message"
// @Router /api/v1/folder/delete_email [delete]
func (h *FolderHandler) DeleteEmailInFolder(w http.ResponseWriter, r *http.Request) {
	var newFolderEmail folderApi.FolderEmail
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newFolderEmail)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	profileId, err := h.Sessions.GetProfileIDBySessionID(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad session")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	folderServiceClient := proto.NewFolderServiceClient(conn)

	CheckFolderProfileDataProto, err := folderServiceClient.CheckFolderProfile(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.FolderProfile{
			FolderID:  newFolderEmail.FolderID,
			ProfileID: profileId,
		},
	)
	if err != nil || !CheckFolderProfileDataProto.Status {
		response.HandleError(w, http.StatusBadRequest, "ProfileID and FolderID not found")
		return
	}

	CheckEmailProfileDataProto, err := folderServiceClient.CheckEmailProfile(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.EmailProfile{
			EmailID:   newFolderEmail.EmailID,
			ProfileID: profileId,
		},
	)
	if err != nil || !CheckEmailProfileDataProto.Status {
		response.HandleError(w, http.StatusBadRequest, "ProfileID and EmailID not found")
		return
	}

	folderDataProto, err := folderServiceClient.DeleteEmailInFolder(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.FolderEmail{
			FolderID: newFolderEmail.FolderID,
			EmailID:  newFolderEmail.EmailID,
		},
	)
	if err != nil || !folderDataProto.Status {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete FolderEmail message")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": folderDataProto.Status})
}

// GetAllEmailsInFolder get all emails in folder.
// @Summary GetAllEmailsInFolder get all emails in folder
// @Description GetAllEmailsInFolder emails in folder users
// @Tags folders
// @Produce json
// @Param id path integer true "ID of the folder"
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "ID of the send folder message"
// @Failure 400 {object} response.Response "Bad JSON in request"
// @Failure 401 {object} response.Response "Not Authorized"
// @Failure 500 {object} response.Response "Failed to get all emails in folder"
// @Router /api/v1/folder/all_emails/{id} [get]
func (h *FolderHandler) GetAllEmailsInFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	profileId, err := h.Sessions.GetProfileIDBySessionID(r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad session")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.FolderService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	folderServiceClient := proto.NewFolderServiceClient(conn)
	CheckFolderProfileDataProto, err := folderServiceClient.CheckFolderProfile(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.FolderProfile{
			FolderID:  uint32(id),
			ProfileID: profileId,
		},
	)
	if err != nil || !CheckFolderProfileDataProto.Status {
		response.HandleError(w, http.StatusBadRequest, "ProfileID and FolderID not found")
		return
	}

	emailsDataProto, err := folderServiceClient.GetAllEmailsInFolder(
		metadata.NewOutgoingContext(r.Context(), metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.GetAllEmailsInFolderData{FolderID: uint32(id), ProfileID: profileId, Limit: 0, Offset: 0},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to get all emails in folder")
		return
	}

	emailsCore := proto_converters.ObjectsEmailConvertProtoInCore(emailsDataProto)

	emailsApi := make([]*folderApi.Email, 0, len(emailsCore))
	for _, email := range emailsCore {
		emailsApi = append(emailsApi, converters.EmailConvertCoreInApi(*email))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"folders": emailsApi})
}
