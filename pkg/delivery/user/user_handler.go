package user

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"io"
	"mail/pkg/delivery"
	"mail/pkg/delivery/converters"
	api "mail/pkg/delivery/models"
	domainSession "mail/pkg/domain/delivery"
	"mail/pkg/domain/usecase"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	UHandler                        = &UserHandler{}
	requestIDContextKey interface{} = "requestid"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	UserUseCase usecase.UserUseCase
	Sessions    domainSession.SessionsManager
}

// InitializationUserHandler initializes the user handler with the provided user handler.
func InitializationUserHandler(userHandler *UserHandler) {
	UHandler = userHandler
}

// sanitizeString sanitizes the provided string using the UGCPolicy from the bluemonday package.
func sanitizeString(str string) string {
	p := bluemonday.UGCPolicy()

	return p.Sanitize(str)
}

// VerifyAuth verifies user authentication.
// @Summary Verify user authentication
// @Description Verify user authentication using sessions
// @Tags users
// @Produce json
// @Success 200 {object} delivery.Response "OK"
// @Failure 401 {object} delivery.Response "Not Authorized"
// @Router /api/v1/verify-auth [get]
func (uh *UserHandler) VerifyAuth(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	sessionUser := uh.Sessions.GetSession(r, requestID)
	w.Header().Set("X-Csrf-Token", sessionUser.CsrfToken)

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "OK"})
}

// GetUserBySession retrieves the user associated with the current session.
// @Summary Get user by session
// @Description Retrieve the user associated with the current session
// @Tags users
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} delivery.Response "User details"
// @Failure 401 {object} delivery.ErrorResponse "Not Authorized"
// @Failure 500 {object} delivery.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/get [get]
func (uh *UserHandler) GetUserBySession(w http.ResponseWriter, r *http.Request) {
	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}
	sessionUser := uh.Sessions.GetSession(r, requestID)
	userData, err := uh.UserUseCase.GetUserByID(sessionUser.UserID, requestID)
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userData.Login = sanitizeString(strings.TrimSpace(userData.Login))
	userData.FirstName = sanitizeString(strings.TrimSpace(userData.FirstName))
	userData.Surname = sanitizeString(strings.TrimSpace(userData.Surname))
	userData.Patronymic = sanitizeString(strings.TrimSpace(userData.Patronymic))
	userData.AvatarID = sanitizeString(strings.TrimSpace(userData.AvatarID))
	userData.PhoneNumber = sanitizeString(strings.TrimSpace(userData.PhoneNumber))
	userData.Description = sanitizeString(strings.TrimSpace(userData.Description))

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"user": converters.UserConvertCoreInApi(*userData)})
}

// UpdateUserData handles requests to update user data.
// @Summary Update user data
// @Description Handles requests to update user data.
// @Tags users
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param updatedUser body delivery.UserSwag true "Updated user data"
// @Success 200 {object} delivery.Response "User data updated successfully"
// @Failure 400 {object} delivery.ErrorResponse "Invalid request body"
// @Failure 401 {object} delivery.ErrorResponse "Not authorized"
// @Failure 500 {object} delivery.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/update [put]
func (uh *UserHandler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	var updatedUser api.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser.Login = sanitizeString(updatedUser.Login)
	updatedUser.FirstName = sanitizeString(updatedUser.FirstName)
	updatedUser.Surname = sanitizeString(updatedUser.Surname)
	updatedUser.Patronymic = sanitizeString(updatedUser.Patronymic)
	updatedUser.AvatarID = sanitizeString(updatedUser.AvatarID)
	updatedUser.PhoneNumber = sanitizeString(updatedUser.PhoneNumber)
	updatedUser.Description = sanitizeString(updatedUser.Description)

	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	sessionUser := uh.Sessions.GetSession(r, requestID)
	if sessionUser.UserID != updatedUser.ID {
		delivery.HandleError(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	userUpdated, err := uh.UserUseCase.UpdateUser(converters.UserConvertApiInCore(updatedUser), requestID)
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"user": converters.UserConvertCoreInApi(*userUpdated)})
}

// DeleteUserData handles requests to delete user data.
// @Summary Delete user data
// @Description Handles requests to delete user data.
// @Tags users
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param id path uint32 true "User ID to delete"
// @Success 200 {object} delivery.Response "User data deleted successfully"
// @Failure 400 {object} delivery.ErrorResponse "Invalid user ID"
// @Failure 401 {object} delivery.ErrorResponse "Not authorized"
// @Failure 500 {object} delivery.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/delete/{id} [delete]
func (uh *UserHandler) DeleteUserData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		delivery.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}

	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	sessionUser := uh.Sessions.GetSession(r, requestID)
	if sessionUser.UserID != uint32(userID) {
		delivery.HandleError(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	deleted, err := uh.UserUseCase.DeleteUserByID(uint32(userID), requestID)
	if err != nil {
		delivery.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if !deleted {
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to delete user data")
		return
	}

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"message": "User data deleted successfully"})
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
// @Success 200 {object} delivery.Response "File uploaded and saved successfully"
// @Failure 400 {object} delivery.ErrorResponse "Error processing file or failed to get file"
// @Failure 500 {object} delivery.ErrorResponse "Internal Server Error"
// @Router /api/v1/user/avatar/upload [post]
func (uh *UserHandler) UploadUserAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(5 * 1024 * 1024)
	if err != nil {
		fmt.Println(err)
		delivery.HandleError(w, http.StatusBadRequest, "Error processing file")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		delivery.HandleError(w, http.StatusBadRequest, "Failed to get file")
		return
	}
	defer file.Close()

	if handler.Size > (5 * 1024 * 1024) {
		fmt.Println(err)
		delivery.HandleError(w, http.StatusInternalServerError, "Failed to get file")
		return
	}

	fileExt := filepath.Ext(handler.Filename)
	uniqueFileName := generateUniqueFileName(fileExt)
	outFile, err := os.Create("./avatars/" + uniqueFileName)
	if err != nil {
		fmt.Println(err)
		delivery.HandleError(w, http.StatusInternalServerError, "Error creating file")
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		fmt.Println(err)
		delivery.HandleError(w, http.StatusInternalServerError, "Error saving file")
		return
	}

	requestID, ok := r.Context().Value(requestIDContextKey).(string)
	if !ok {
		requestID = "none"
	}

	sessionUser := uh.Sessions.GetSession(r, requestID)
	userData, err := uh.UserUseCase.GetUserByID(sessionUser.UserID, requestID)
	if err != nil {
		fmt.Println(err)
		delivery.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userData.AvatarID = "http://mailhub.su:8080/media/" + uniqueFileName
	userUpdated, err := uh.UserUseCase.UpdateUser(userData, requestID)
	if err != nil {
		fmt.Println(err)
		delivery.HandleError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	fmt.Println(userUpdated)

	delivery.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "File is uploaded and saved"})
}

// generateUniqueFileName generates a unique file name based on the current time, random number, and specified format.
func generateUniqueFileName(format string) string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(1000)

	currentTime := time.Now().Format("20060102_150405")
	uniqueFileName := fmt.Sprintf("%s_%d%s", currentTime, randomNum, format)

	return uniqueFileName
}
