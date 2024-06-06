package http

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"mail/internal/microservice/models/domain_models"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"
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
	OAHandler           = &OAuthHandler{}
	AUTH_URL            = ""
	APP_ID              = ""
	APP_KEY             = ""
	API_URL             = ""
	REDIRECT_URL_SIGNUP = ""
	REDIRECT_URL_LOGIN  = ""
	mapVKIDToken        = make(map[uint32]string)
)

type Response struct {
	Response []struct {
		VKId      int    `json:"id"`
		Name      string `json:"first_name"`
		LastName  string `json:"last_name"`
		Sex       int    `json:"sex"`
		BirthDate string `json:"bdate"`
	}
}

// OAuthHandler handles user-related HTTP requests.
type OAuthHandler struct {
	Sessions          domainSession.SessionsManager
	UserServiceClient user_proto.UserServiceClient
}

// InitializationOAuthHandler initializes the user handler with the provided user handler.
func InitializationOAuthHandler(oauthHandler *OAuthHandler) {
	OAHandler = oauthHandler
}

// GetSignUpURLVK url auth vk.
// @Summary URL VK
// @Description Handles user signup.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Successful"
// @Failure 500 {object} response.ErrorResponse "Failed to get url"
// @Router /api/v1/testAuth/auth-vk/getAuthUrlSignUpVK [get]
func (ah *OAuthHandler) GetSignUpURLVK(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf(AUTH_URL, APP_ID, REDIRECT_URL_SIGNUP)
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"AuthURL": url})
}

// GetLoginURLVK url auth vk.
// @Summary URL VK
// @Description Handles user login.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Successful"
// @Failure 500 {object} response.ErrorResponse "Failed to get url"
// @Router /api/v1/testAuth/auth-vk/getAuthUrlLoginVK [get]
func (ah *OAuthHandler) GetLoginURLVK(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf(AUTH_URL, APP_ID, REDIRECT_URL_LOGIN)
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"AuthURL": url})
}

// AuthVK handles user auth VK.
// @Summary User auth VK
// @Description Handles user auth VK.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Param code path string true "Code of the oauth message"
// @Success 200 {object} response.Response "Auth successful"
// @Failure 500 {object} response.ErrorResponse "Failed to auth user"
// @Router /api/v1/testAuth/auth-vk/auth/{code} [get]
func (ah *OAuthHandler) AuthVK(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	fmt.Println("Code: ", code)
	if !ok {
		response.HandleError(w, http.StatusBadRequest, "Bad id in request")
		return
	}
	ctx := r.Context()

	conf := GetConfOauth2(REDIRECT_URL_SIGNUP)

	// vk_mock
	if code == "855ab871bba885204e" {
		vkUser := &api.VKUser{
			FirstName: "Max",
			Surname:   "Frelih",
			Gender:    domain_models.GetGenderTypeInt(2),
			VKId:      uint32(1234567),
		}
		randToken := make([]byte, 16)
		_, err := rand.Read(randToken)
		if err != nil {
			fmt.Println("Error reading random numbers:", err)
			return
		}
		authToken := fmt.Sprintf("%x", randToken)
		mapVKIDToken[vkUser.VKId] = authToken
		w.Header().Set("AuthToken", authToken)
		response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"VKUser": vkUser})
		return
	}

	if code == "" {
		response.HandleError(w, http.StatusBadRequest, "wrong code")
		return
	}

	vkUser, status, err := GetDataUser(*conf, code, ctx)
	if err != nil {
		response.HandleError(w, status, "failed get user data")
		return
	}

	fmt.Println(vkUser.Birthday)

	randToken := make([]byte, 16)
	_, err = rand.Read(randToken)
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "failed to generate random token")
		return
	}
	authToken := fmt.Sprintf("%x", randToken)
	mapVKIDToken[vkUser.VKId] = authToken
	w.Header().Set("AuthToken", authToken)

	fmt.Println("authToken: ", authToken)
	fmt.Println("name: ", vkUser.FirstName, ", surname: ", vkUser.Surname, ", gender: ", vkUser.Gender, ", vkId: ", vkUser.VKId)
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"VKUser": vkUser})
}

// SignupVK handles user signup VK.
// @Summary User signup VK
// @Description Handles user signup VK.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Param AuthToken header string true "Auth Token"
// @Param newUser body response.UserVKSwag true "New user details for signup"
// @Success 200 {object} response.Response "Signup successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to add user"
// @Router /api/v1/testAuth/auth-vk/signupVK [post]
func (ah *OAuthHandler) SignupVK(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newUser api.VKUser
	if err := newUser.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	authToken := r.Header.Get("AuthToken")
	if authToken != mapVKIDToken[newUser.VKId] {
		response.HandleError(w, http.StatusBadRequest, "failed authToken")
		return
	}

	newUser.FirstName = sanitize.SanitizeString(newUser.FirstName)
	newUser.Surname = sanitize.SanitizeString(newUser.Surname)
	newUser.Login = sanitize.SanitizeString(newUser.Login)

	if validUtil.IsEmpty(newUser.Login) || validUtil.IsEmpty(newUser.FirstName) || validUtil.IsEmpty(newUser.Surname) || !domain.IsValidGender(newUser.Gender) {
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

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.AuthService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()

	authServiceClient := auth_proto.NewAuthServiceClient(conn)
	_, errStatus := authServiceClient.SignupVK(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.SignupVKRequest{
			Login:     newUser.Login,
			Firstname: newUser.FirstName,
			Surname:   newUser.Surname,
			Birthday:  timestamppb.New(newUser.Birthday),
			Gender:    domain.GetGender(newUser.Gender),
			VkId:      newUser.VKId,
		},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to add user")
		return
	}

	sessionId, errStatus := authServiceClient.LoginVK(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.LoginVKRequest{VkId: newUser.VKId},
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

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})
}

// LoginVK handles user login.
// @Summary LoginVK User
// @Description LoginVK Handles user.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Param code path string true "Code of the oauth message"
// @Success 200 {object} response.Response "Login successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Failed to create session"
// @Router /api/v1/testAuth/auth-vk/loginVK/{code} [get]
func (ah *OAuthHandler) LoginVK(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	fmt.Println("Code: ", code)
	if !ok || code == "" {
		response.HandleError(w, http.StatusBadRequest, "Bad code in request")
		return
	}
	ctx := r.Context()
	conf := GetConfOauth2(REDIRECT_URL_LOGIN)

	var userVK *api.VKUser
	if code == "855ab871bba885204e" {
		userVK = &api.VKUser{
			FirstName: "Max",
			Surname:   "Frelih",
			Gender:    domain_models.GetGenderTypeInt(2),
			VKId:      1234567,
		}
	} else {
		userVk, status, err := GetDataUser(*conf, code, ctx)
		if err != nil {
			response.HandleError(w, status, "failed get user data")
			return
		}
		userVK = userVk
	}

	if userVK.VKId <= 0 {
		response.HandleError(w, http.StatusBadRequest, "bad VKId")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.AuthService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()

	authServiceClient := auth_proto.NewAuthServiceClient(conn)
	sessionId, errStatus := authServiceClient.LoginVK(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.LoginVKRequest{VkId: userVK.VKId},
	)
	if errStatus != nil {
		randToken := make([]byte, 16)
		_, err = rand.Read(randToken)
		if err != nil {
			response.HandleError(w, http.StatusInternalServerError, "failed to generate random token")
			return
		}
		authToken := fmt.Sprintf("%x", randToken)
		mapVKIDToken[userVK.VKId] = authToken
		w.Header().Set("AuthToken", authToken)
		response.HandleSuccess(w, http.StatusUnauthorized, map[string]interface{}{"VKUser": userVK})
		return
	}

	er := ah.Sessions.SetSession(sessionId.SessionId, w, r, r.Context())
	if er != nil {
		response.HandleError(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Login successful"})
}

func GetConfOauth2(redirectUrl string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     APP_ID,
		ClientSecret: APP_KEY,
		RedirectURL:  redirectUrl,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.vk.com/authorize",
			TokenURL: "https://oauth.vk.com/access_token",
		},
		Scopes: []string{"email"},
	}
}

func GetDataUser(conf oauth2.Config, code string, ctx context.Context) (*api.VKUser, int, error) {
	token, err := conf.Exchange(ctx, code)
	fmt.Println("Token: ", token)
	if err != nil {
		fmt.Println(err)
		return &api.VKUser{}, 400, fmt.Errorf("cannot exchange")
	}

	client := conf.Client(ctx, token)
	resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
	if err != nil {
		fmt.Println("cannot request data")
		return &api.VKUser{}, 400, fmt.Errorf("cannot request data")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("cannot read buffer")
		return &api.VKUser{}, 500, fmt.Errorf("cannot read buffer")
	}

	data := &Response{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return &api.VKUser{}, 400, fmt.Errorf("cannot unmarshal response")
	}

	fmt.Println("Data: ", data.Response[0].BirthDate)

	date, err := time.Parse("2.1.2006", data.Response[0].BirthDate)
	if err != nil {
		fmt.Println(err)
		return &api.VKUser{}, 400, fmt.Errorf("failed to parse date")
	}

	vkUser := &api.VKUser{
		FirstName: data.Response[0].Name,
		Surname:   data.Response[0].LastName,
		Gender:    domain_models.GetGenderTypeInt(data.Response[0].Sex),
		Birthday:  date,
		VKId:      uint32(data.Response[0].VKId),
	}

	return vkUser, 200, nil
}
