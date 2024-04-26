package http

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc/metadata"
	"mail/internal/microservice/folder/proto"
	"mail/internal/microservice/models/proto_converters"
	converters "mail/internal/models/delivery_converters"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"

	folderUsecase "mail/internal/microservice/folder/interface"
	folderApi "mail/internal/models/delivery_models"
	"mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"net/http"
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
	fmt.Println(folderData)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"folder": converters.FolderConvertCoreInApi(*folderData)})
	return
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
	return
}
