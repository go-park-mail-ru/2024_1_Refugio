package http

import (
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"

	"mail/internal/pkg/utils/sanitize"

	auth_proto "mail/internal/microservice/auth/proto"
	domain "mail/internal/microservice/models/domain_models"
	user_proto "mail/internal/microservice/user/proto"
	api "mail/internal/models/delivery_models"
	response "mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	validUtil "mail/internal/pkg/utils/validators"
)

var (
	AHandler = &AuthHandler{}
)

// AuthHandler handles user-related HTTP requests.
type AuthHandler struct {
	Sessions          domainSession.SessionsManager
	AuthServiceClient auth_proto.AuthServiceClient
	UserServiceClient user_proto.UserServiceClient
}

// Login handles user login.
// @Summary Login User
// @Description Login Handles user.
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body response.UserSwag true "User credentials for login"
// @Success 200 {object} response.Response "Login successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Failed to create session"
// @Router /api/v1/auth/login [post]
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var credentials api.User
	if err := credentials.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	credentials.Login = sanitize.SanitizeString(credentials.Login)
	credentials.Password = sanitize.SanitizeString(credentials.Password)

	if validUtil.IsEmpty(credentials.Login) || validUtil.IsEmpty(credentials.Password) {
		response.HandleError(w, http.StatusInternalServerError, "All fields must be filled in")
		return
	}

	if !validUtil.IsValidEmailFormat(credentials.Login) {
		response.HandleError(w, http.StatusBadRequest, "Domain in the login is not suitable")
		return
	}

	sessionId, errStatus := ah.AuthServiceClient.Login(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.LoginRequest{Login: credentials.Login, Password: credentials.Password},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusUnauthorized, "Login failed")
		return
	}

	er := ah.Sessions.SetSession(sessionId.SessionId, w, r, r.Context())
	if er != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Login successful"})
}

// Signup handles user signup.
// @Summary User signup
// @Description Handles user signup.
// @Tags auth
// @Accept json
// @Produce json
// @Param newUser body response.UserSwag true "New user details for signup"
// @Success 200 {object} response.Response "Signup successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to add user"
// @Router /api/v1/auth/signup [post]
func (ah *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newUser api.User
	if err := newUser.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newUser.Login = sanitize.SanitizeString(newUser.Login)
	newUser.Password = sanitize.SanitizeString(newUser.Password)
	newUser.FirstName = sanitize.SanitizeString(newUser.FirstName)
	newUser.Surname = sanitize.SanitizeString(newUser.Surname)
	newUser.Patronymic = sanitize.SanitizeString(newUser.Patronymic)
	newUser.PhoneNumber = sanitize.SanitizeString(newUser.PhoneNumber)
	newUser.Description = sanitize.SanitizeString(newUser.Description)
	newUser.AvatarID = sanitize.SanitizeString(newUser.AvatarID)

	if validUtil.IsEmpty(newUser.Login) || validUtil.IsEmpty(newUser.Password) || validUtil.IsEmpty(newUser.FirstName) || validUtil.IsEmpty(newUser.Surname) || !domain.IsValidGender(newUser.Gender) {
		response.HandleError(w, http.StatusBadRequest, "All fields must be filled in")
		return
	}

	if !validUtil.IsValidEmailFormat(newUser.Login) {
		response.HandleError(w, http.StatusBadRequest, "Domain in the login is not suitable")
		return
	}

	userExists, err := ah.UserServiceClient.GetUserByOnlyLogin(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&user_proto.GetUserByOnlyLoginRequest{
			Login: newUser.Login,
		},
	)
	if userExists != nil || err == nil {
		response.HandleError(w, http.StatusBadRequest, "User already exists")
		return
	}

	_, errStatus := ah.AuthServiceClient.Signup(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.SignupRequest{Login: newUser.Login,
			Password:    newUser.Password,
			Firstname:   newUser.FirstName,
			Surname:     newUser.Surname,
			Patronymic:  newUser.Patronymic,
			Birthday:    timestamppb.New(newUser.Birthday),
			Gender:      domain.GetGender(newUser.Gender),
			Avatar:      newUser.AvatarID,
			Description: newUser.Description,
			PhoneNumber: newUser.PhoneNumber,
		},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to add user")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})
}

// Logout handles user logout.
// @Summary User logout
// @Description Handles user logout.
// @Tags auth
// @Produce json
// @Success 200 {object} response.Response "Logout successful"
// @Router /api/v1/auth/logout [post]
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	err := ah.Sessions.DestroyCurrent(w, r, r.Context())
	if err != nil {
		response.HandleError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Logout successful"})
}
