package http

import (
	"fmt"
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"

	"mail/cmd/configs"
	"mail/internal/microservice/models/proto_converters"
	"mail/internal/microservice/user/proto"
	"mail/internal/pkg/utils/check_image"
	"mail/internal/pkg/utils/generate_filename"
	"mail/internal/pkg/utils/sanitize"

	user_proto "mail/internal/microservice/user/proto"
	converters "mail/internal/models/delivery_converters"
	api "mail/internal/models/delivery_models"
	response "mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
)

var (
	UHandler                        = &UserHandler{}
	requestIDContextKey interface{} = "requestid"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	Sessions          domainSession.SessionsManager
	UserServiceClient user_proto.UserServiceClient
	MinioClient       *minio.Client
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

	userDataProto, err := uh.UserServiceClient.GetUser(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.GetUserRequest{Id: sessionUser.UserID},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	userData := proto_converters.UserConvertProtoInCore(userDataProto.User)

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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var updatedUser api.User
	if err := updatedUser.UnmarshalJSON(body); err != nil {
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

	userUpdateProto, err := uh.UserServiceClient.UpdateUser(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.UpdateUserRequest{User: proto_converters.UserConvertCoreInProto(converters.UserConvertApiInCore(updatedUser))},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	userUpdated := proto_converters.UserConvertProtoInCore(userUpdateProto.User)

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

	deleted, err := uh.UserServiceClient.DeleteUserById(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.DeleteUserByIdRequest{Id: sessionUser.UserID},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if !deleted.Status {
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
	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
		response.HandleError(w, http.StatusBadRequest, "File type not supported")
		return
	}

	tempFile, _, errImg := r.FormFile("file")
	if errImg != nil {
		response.HandleError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer tempFile.Close()
	if !check_image.IsImage(tempFile) {
		response.HandleError(w, http.StatusBadRequest, "File is not an image")
		return
	}

	uniqueFileName := generate_filename.GenerateUniqueFileName(fileExt)
	_, err = uh.MinioClient.PutObject(r.Context(), "photos", uniqueFileName, file, -1, minio.PutObjectOptions{ContentType: handler.Header.Get("Content-Type")})
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Error uploading file to MinIO")
		return
	}

	sessionUser := uh.Sessions.GetSession(r, r.Context())

	avatarURL := fmt.Sprintf(configs.PROTOCOL+"mailhub.su"+"/photos/%s", uniqueFileName)
	_, errAvatar := uh.UserServiceClient.UploadUserAvatar(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.UploadUserAvatarRequest{Id: sessionUser.UserID, Avatar: avatarURL},
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

	deleted, err := uh.UserServiceClient.DeleteUserAvatar(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.DeleteUserAvatarRequest{Id: sessionUser.UserID},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if !deleted.Status {
		response.HandleError(w, http.StatusInternalServerError, "Failed to delete user avatar")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"message": "User avatar deleted successfully"})
}

// GetCountUsers handles requests to get count users.
// @Summary Get count user
// @Description Handles requests to get count user.
// @Tags users
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 401 {object} response.ErrorResponse "Not authorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/count [get]
func (uh *UserHandler) GetCountUsers(w http.ResponseWriter, r *http.Request) {
	usersProto, err := uh.UserServiceClient.GetUsers(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&proto.GetUsersRequest{},
	)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"count": len(usersProto.Users)})
}
