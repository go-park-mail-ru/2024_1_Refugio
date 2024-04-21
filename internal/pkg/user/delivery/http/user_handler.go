package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"mail/internal/microservice/models/proto_converters"
	"mail/internal/microservice/user/interface"
	"mail/internal/microservice/user/proto"
	converters "mail/internal/models/delivery_converters"
	api "mail/internal/models/delivery_models"
	"mail/internal/models/microservice_ports"
	response "mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"mail/internal/pkg/utils/connect_microservice"
	"mail/internal/pkg/utils/generate_filename"
	"mail/internal/pkg/utils/sanitize"
)

var (
	UHandler                        = &UserHandler{}
	requestIDContextKey interface{} = "requestid"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	UserUseCase  _interface.UserUseCase
	Sessions     domainSession.SessionsManager
	MicroService proto.UserServiceClient
}

// InitializationUserHandler initializes the user handler with the provided user handler.
func InitializationUserHandler(userHandler *UserHandler) {
	UHandler = userHandler
}

// VerifyAuth verifies user authentication.
// @Summary Verify user authentication
// @Description Verify user authentication using sessions
// @Tags users
// @Produce json
// @Success 200 {object} response.Response "OK"
// @Failure 401 {object} response.Response "Not Authorized"
// @Router /api/v1/verify-auth [get]
func (uh *UserHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	sessionUser := uh.Sessions.GetSession(r, r.Context())
	w.Header().Set("X-Csrf-Token", sessionUser.CsrfToken)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "OK"})
}

// GetUserBySession retrieves the user associated with the current session.
// @Summary Get user by session
// @Description Retrieve the user associated with the current session
// @Tags users
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "User details"
// @Failure 401 {object} response.ErrorResponse "Not Authorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/get [get]
func (uh *UserHandler) GetUserBySession(w http.ResponseWriter, r *http.Request) {
	sessionUser := uh.Sessions.GetSession(r, r.Context())

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	userServiceClient := proto.NewUserServiceClient(conn)
	userDataProto, err := userServiceClient.GetUser(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.UserId{Id: sessionUser.UserID},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	userData := proto_converters.UserConvertProtoInCore(*userDataProto)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"user": converters.UserConvertCoreInApi(*userData)})
}

// UpdateUserData handles requests to update user data.
// @Summary Update user data
// @Description Handles requests to update user data.
// @Tags users
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param updatedUser body response.UserSwag true "Updated user data"
// @Success 200 {object} response.Response "User data updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Not authorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/update [put]
func (uh *UserHandler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	var updatedUser api.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser.Login = sanitize.SanitizeString(updatedUser.Login)
	updatedUser.FirstName = sanitize.SanitizeString(updatedUser.FirstName)
	updatedUser.Surname = sanitize.SanitizeString(updatedUser.Surname)
	updatedUser.Patronymic = sanitize.SanitizeString(updatedUser.Patronymic)
	updatedUser.AvatarID = sanitize.SanitizeString(updatedUser.AvatarID)
	updatedUser.PhoneNumber = sanitize.SanitizeString(updatedUser.PhoneNumber)
	updatedUser.Description = sanitize.SanitizeString(updatedUser.Description)

	sessionUser := uh.Sessions.GetSession(r, r.Context())
	if sessionUser.UserID != updatedUser.ID {
		response.HandleError(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	userServiceClient := proto.NewUserServiceClient(conn)
	userUpdateProto, err := userServiceClient.UpdateUser(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		proto_converters.UserConvertCoreInProto(*converters.UserConvertApiInCore(updatedUser)),
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	userUpdated := proto_converters.UserConvertProtoInCore(*userUpdateProto)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"user": converters.UserConvertCoreInApi(*userUpdated)})
}

// DeleteUserData handles requests to delete user data.
// @Summary Delete user data
// @Description Handles requests to delete user data.
// @Tags users
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path uint32 true "User ID to delete"
// @Success 200 {object} response.Response "User data deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 401 {object} response.ErrorResponse "Not authorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/delete/{id} [delete]
func (uh *UserHandler) DeleteUserData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	sessionUser := uh.Sessions.GetSession(r, r.Context())
	if sessionUser.UserID != uint32(userID) {
		response.HandleError(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	userServiceClient := proto.NewUserServiceClient(conn)
	deleted, err := userServiceClient.DeleteUserById(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.UserId{Id: sessionUser.UserID},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if !deleted.DeleteStatus {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete user data")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"message": "User data deleted successfully"})
}

// UploadUserAvatar handles requests to upload user avatar.
// @Summary Upload user avatar
// @Description Handles requests to upload user avatar.
// @Tags users
// @Accept multipart/form-data
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param file formData file true "Avatar file to upload"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "File uploaded and saved successfully"
// @Failure 400 {object} response.ErrorResponse "Error processing file or failed to get file"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/avatar/upload [post]
func (uh *UserHandler) UploadUserAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 * 1024 * 1024)
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

	if handler.Size > (5 * 1024 * 1024) {
		response.HandleError(w, http.StatusInternalServerError, "Failed to get file")
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	uniqueFileName := generate_filename.GenerateUniqueFileName(fileExt)
	outFile, err := os.Create("./avatars/" + uniqueFileName)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error creating file")
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error saving file")
		return
	}

	sessionUser := uh.Sessions.GetSession(r, r.Context())

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	userServiceClient := proto.NewUserServiceClient(conn)
	userDataProto, err := userServiceClient.GetUser(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.UserId{Id: sessionUser.UserID},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userDataProto.Avatar = "https://mailhub.su/media/" + uniqueFileName
	_, errAvatar := userServiceClient.UploadUserAvatar(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.UserAvatar{Id: userDataProto.Id, Avatar: userDataProto.Avatar},
	)
	if errAvatar != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "File is uploaded and saved"})
}

// DeleteUserAvatar handles requests to delete user avatar.
// @Summary Delete user avatar
// @Description Handles requests to delete user avatar.
// @Tags users
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "User avatar deleted successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 401 {object} response.ErrorResponse "Not authorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/avatar/delete [delete]
func (uh *UserHandler) DeleteUserAvatar(w http.ResponseWriter, r *http.Request) {
	sessionUser := uh.Sessions.GetSession(r, r.Context())

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.UserService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	userServiceClient := proto.NewUserServiceClient(conn)
	deleted, err := userServiceClient.DeleteUserAvatar(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.UserId{Id: sessionUser.UserID},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if !deleted.DeleteStatus {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete user avatar")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"message": "User avatar deleted successfully"})
}
