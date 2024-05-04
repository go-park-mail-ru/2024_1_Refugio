package http

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	auth_proto "mail/internal/microservice/auth/proto"
	api "mail/internal/models/delivery_models"
	"mail/internal/models/microservice_ports"
	"mail/internal/pkg/utils/connect_microservice"
	"mail/internal/pkg/utils/sanitize"
	validUtil "mail/internal/pkg/utils/validators"
	"math/rand"
	"net/http"
	"time"

	"mail/internal/microservice/models/domain_models"
	domain "mail/internal/microservice/models/domain_models"
	response "mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
)

var (
	OAHandler                       = &OAuthHandler{}
	requestIDContextKey interface{} = "requestid"
	AUTH_URL                        = "https://oauth.vk.com/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=email"
	APP_ID                          = "51916655"
	APP_KEY                         = "oz3r7Pyakfeg25JpJsQV"
	API_URL                         = "https://api.vk.com/method/users.get?fields=id,photo_max,email,sex,bdate&access_token=%s&v=5.131"
	REDIRECT_URL_SIGNUP             = "https://mailhub.su/testAuth/auth-vk/auth"
	REDIRECT_URL_LOGIN              = "https://mailhub.su/testAuth/auth-vk/loginVK"
	mepVKIDToken                    = make(map[uint32]string)
)

type Response struct {
	Response []struct {
		VKId      int    `json:"id"`
		Name      string `json:"first_name"`
		LastName  string `json:"last_name"`
		Sex       int    `json:"sex"`
		BirthDate string `json:"bdate"`
		//Photo     string `json:"photo_max"`
	}
}

// AuthHandler handles user-related HTTP requests.
type OAuthHandler struct {
	Sessions domainSession.SessionsManager
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
// @Success 200 {object} response.Response "Auth successful"
// @Failure 500 {object} response.ErrorResponse "Failed to auth user"
// @Router /api/v1/testAuth/auth-vk/auth [get]
func (ah *OAuthHandler) AuthVK(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthVK")
	ctx := r.Context()
	code := r.FormValue("code")
	conf := GetConfOauth2(REDIRECT_URL_SIGNUP)

	if code == "" {
		response.HandleError(w, http.StatusBadRequest, "wrong code")
		return
	}

	vkUser, status, err := GetDataUser(*conf, code, ctx)
	if err != nil {
		response.HandleError(w, status, "failed get user data")
		return
	}

	randToken := make([]byte, 16)
	rand.Read(randToken)
	authToken := fmt.Sprintf("%x", randToken)
	mepVKIDToken[vkUser.VKId] = authToken
	w.Header().Set("AuthToken", authToken)

	fmt.Println("name: ", vkUser.FirstName, ", surname: ", vkUser.Surname, ", gender: ", vkUser.Gender, ", vkId: ", vkUser.VKId)
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"VKUser": vkUser})
}

// SignupVK handles user signup VK.
// @Summary User signup VK
// @Description Handles user signup VK.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Param Auth-Token header string true "Auth Token"
// @Param newUser body response.UserVKSwag true "New user details for signup"
// @Success 200 {object} response.Response "Signup successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to add user"
// @Router /api/v1/testAuth/auth-vk/signupVK [post]
func (ah *OAuthHandler) SignupVK(w http.ResponseWriter, r *http.Request) {
	var newUser api.VKUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	//mepVKIDToken[123] = "123"

	authToken := r.Header.Get("Auth-Token")
	if authToken != mepVKIDToken[newUser.VKId] {
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

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": "Signup successful"})
}

// LoginVK handles user login.
// @Summary LoginVK User
// @Description LoginVK Handles user.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Param code query string true "code from oauth"
// @Success 200 {object} response.Response "Login successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Failed to create session"
// @Router /api/v1/testAuth/auth-vk/loginVK [get]
func (ah *OAuthHandler) LoginVK(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginVK")
	ctx := r.Context()
	code := r.FormValue("code")
	conf := GetConfOauth2(REDIRECT_URL_LOGIN)
	if code == "" {
		response.HandleError(w, http.StatusBadRequest, "wrong code")
		return
	}

	userVK, status, err := GetDataUser(*conf, code, ctx)
	if err != nil {
		response.HandleError(w, status, "failed get user data")
		return
	}

	/*
		userVK := &api.VKUser{
			VKId: 123,
		}
	*/
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
		return &api.VKUser{}, 400, fmt.Errorf("cannot exchange")
	}

	fmt.Println("TOKEN OK")

	client := conf.Client(ctx, token)
	resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
	if err != nil {
		return &api.VKUser{}, 400, fmt.Errorf("cannot request data")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &api.VKUser{}, 500, fmt.Errorf("cannot read buffer")
	}

	data := &Response{}
	json.Unmarshal(body, data)

	birthdayTime, err := time.Parse("2006-01-02", data.Response[0].BirthDate)
	if err != nil {
		return &api.VKUser{}, 500, fmt.Errorf("bad BirthDate")
	}

	vkUser := &api.VKUser{
		FirstName: data.Response[0].Name,
		Surname:   data.Response[0].LastName,
		Gender:    domain_models.GetGenderTypeInt(data.Response[0].Sex),
		Birthday:  birthdayTime,
		VKId:      uint32(data.Response[0].VKId),
	}

	return vkUser, 200, nil
}
