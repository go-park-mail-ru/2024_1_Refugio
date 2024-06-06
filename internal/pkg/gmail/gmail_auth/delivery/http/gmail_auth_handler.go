package http

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"os"

	"mail/internal/models/response"
	"mail/internal/pkg/utils/sanitize"

	auth_proto "mail/internal/microservice/auth/proto"
	domain "mail/internal/microservice/models/domain_models"
	user_proto "mail/internal/microservice/user/proto"
	api "mail/internal/models/delivery_models"
	domainSession "mail/internal/pkg/session/interface"
	validUtil "mail/internal/pkg/utils/validators"
)

var (
	MapOAuthCongig                  = make(map[string]*gmail.Service)
	requestIDContextKey interface{} = "requestid"
)

// GMailAuthHandler handles user-related HTTP requests.
type GMailAuthHandler struct {
	Sessions          domainSession.SessionsManager
	AuthServiceClient auth_proto.AuthServiceClient
	UserServiceClient user_proto.UserServiceClient
} // struct

// GoogleAuth handles user auth.
// @Summary GoogleAuth User
// @Description GoogleAuth Handles user.
// @Tags auth-gmail
// @Accept json
// @Produce json
// @Param code query string true "code from oauth"
// @Success 200 {object} response.Response "Auth successful"
// @Failure 404 {object} response.Response "User not fount"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Failed to create session"
// @Router /api/v1/auth/gAuth [get]
func (g *GMailAuthHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")

	b, err := os.ReadFile("cmd/configs/credentials_localhost.json")
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to read client secret file")
		return
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to parse client secret file to config")
		return
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve token from web")
		return
	}

	ctx := context.Background()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve Gmail client")
		return
	}

	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to retrieve profile user")
		return
	}

	tokFile := "token.json"
	err = saveToken(tokFile, tok)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to cache oauth token")
		return
	}

	MapOAuthCongig[profile.EmailAddress] = srv

	userDataProto, err := g.UserServiceClient.GetUserByOnlyLogin(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&user_proto.GetUserByOnlyLoginRequest{Login: profile.EmailAddress},
	)
	if userDataProto == nil || err != nil {
		response.HandleSuccess(w, http.StatusNotFound, map[string]interface{}{"Status": "UserNotFound", "Login": profile.EmailAddress})
		return
	}

	sessionId, errStatus := g.AuthServiceClient.LoginOtherMail(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.LoginOtherMailRequest{
			Id: userDataProto.User.Id,
		},
	)

	if errStatus != nil {
		response.HandleError(w, http.StatusUnauthorized, "Login failed")
		return
	}

	er := g.Sessions.SetSession(sessionId.SessionId, w, r, r.Context())
	if er != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Status": "OK", "User": userDataProto})
}

// SugnupGMail handles user signup GMail.
// @Summary User signup Gmail
// @Description Handles user signup Gmail.
// @Tags auth-gmail
// @Accept json
// @Produce json
// @Param newUser body response.UserGoogleSwag true "New user details for signup"
// @Success 200 {object} response.Response "Signup successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to add user"
// @Router /api/v1/auth/signupGMailUser [post]
func (g *GMailAuthHandler) SugnupGMail(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newUser api.OtherUser
	if err := newUser.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newUser.Login = sanitize.SanitizeString(newUser.Login)
	newUser.FirstName = sanitize.SanitizeString(newUser.FirstName)
	newUser.Surname = sanitize.SanitizeString(newUser.Surname)
	newUser.Patronymic = sanitize.SanitizeString(newUser.Patronymic)
	newUser.PhoneNumber = sanitize.SanitizeString(newUser.PhoneNumber)
	newUser.Description = sanitize.SanitizeString(newUser.Description)
	newUser.AvatarID = sanitize.SanitizeString(newUser.AvatarID)

	if validUtil.IsEmpty(newUser.Login) || validUtil.IsEmpty(newUser.FirstName) || validUtil.IsEmpty(newUser.Surname) {
		response.HandleError(w, http.StatusBadRequest, "All fields must be filled in")
		return
	}

	_, errStatus := g.AuthServiceClient.SignupOtherMail(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.SignupOtherMailRequest{
			Login:       newUser.Login,
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

	userDataProto, err := g.UserServiceClient.GetUserByOnlyLogin(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&user_proto.GetUserByOnlyLoginRequest{Login: newUser.Login},
	)
	if userDataProto == nil || err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to get user by only login")
		return
	}

	sessionId, errStatus := g.AuthServiceClient.LoginOtherMail(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.LoginOtherMailRequest{
			Id: userDataProto.User.Id,
		},
	)

	if errStatus != nil {
		response.HandleError(w, http.StatusUnauthorized, "Login failed")
		return
	}

	er := g.Sessions.SetSession(sessionId.SessionId, w, r, r.Context())
	if er != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})
}

// GetAuthURL handles get url gmail auth.
// @Summary GetAuthURL
// @Description GetAuthURL Handles url.
// @Tags auth-gmail
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Auth successful"
// @Failure 500 {object} response.ErrorResponse "Failed to create url"
// @Router /api/v1/auth/getAuthURL [get]
func (g *GMailAuthHandler) GetAuthURL(w http.ResponseWriter, r *http.Request) {
	// localhost:8080:
	// b, err := os.ReadFile("cmd/configs/credentials_localhost.json")

	// deploy:
	b, err := os.ReadFile("cmd/configs/credentials_deploy.json")
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to read client secret file")
		return
	}

	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope) // google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "Unable to parse client secret file to config")
		return
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"AuthURL": authURL})
}

func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return fmt.Errorf("Unable to cache oauth token: %v", err)
	}

	return nil
}
