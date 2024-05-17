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
	"log"
	auth_proto "mail/internal/microservice/auth/proto"
	domain "mail/internal/microservice/models/domain_models"
	user_proto "mail/internal/microservice/user/proto"
	api "mail/internal/models/delivery_models"
	"mail/internal/models/microservice_ports"
	"mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"mail/internal/pkg/utils/connect_microservice"
	"mail/internal/pkg/utils/sanitize"
	validUtil "mail/internal/pkg/utils/validators"
	"net/http"
	"os"
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
}

// GoogleAuth handles user auth.
// @Summary GoogleAuth User
// @Description GoogleAuth Handles user.
// @Tags auth-gmail
// @Accept json
// @Produce json
// @Param code query string true "code from oauth"
// @Success 200 {object} response.Response "Auth successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Failed to create session"
// @Router /api/v1/auth/gAuth [get]
func (g *GMailAuthHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	authCode := r.URL.Query().Get("code")

	b, err := os.ReadFile("cmd/configs/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}

	ctx := context.Background()
	client := config.Client(context.Background(), tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	profile, err := srv.Users.GetProfile("me").Do()
	if err != nil {
		log.Fatal(err)
	}

	userDataProto, err := g.UserServiceClient.GetUserByOnlyLogin(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value(requestIDContextKey).(string)})),
		&user_proto.GetUserByOnlyLoginRequest{Login: profile.EmailAddress},
	)
	if userDataProto == nil || err != nil {
		response.HandleError(w, http.StatusBadRequest, "User not found")
		return
	}

	connAuth, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.AuthService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer connAuth.Close()

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

	tokFile := "token.json"
	saveToken(tokFile, tok)
	MapOAuthCongig[profile.EmailAddress] = srv
	fmt.Println(MapOAuthCongig)
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
	var newUser api.OtherUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
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

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})

	/*
		authServiceClient := auth_proto.NewAuthServiceClient(conn)
		sessionId, errStatus := authServiceClient.LoginVK(
			metadata.NewOutgoingContext(r.Context(),
				metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
			&auth_proto.LoginVKRequest{VkId: userVK.VKId},
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

		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Login successful"})
	*/
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
	b, err := os.ReadFile("cmd/configs/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope) // google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Println(authURL)
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"AuthURL": authURL})
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
