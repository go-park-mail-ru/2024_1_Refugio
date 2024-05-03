package http

import (
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
	REDIRECT_URL                    = "https://mailhub.su/auth/auth-vk"
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

// GetAuthURLVK url auth vk.
// @Summary URL VK
// @Description Handles user signup.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Successful"
// @Failure 500 {object} response.ErrorResponse "Failed to get url"
// @Router /api/v1/auth/getAuthUrlVK [get]
func (ah *OAuthHandler) GetAuthURLVK(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf(AUTH_URL, APP_ID, REDIRECT_URL)
	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"AuthURL": url})
}

func AuthVK(w http.ResponseWriter, r *http.Request) (*api.VKUser, int, error) {
	fmt.Println("AuthVK")
	ctx := r.Context()
	code := r.FormValue("code")
	conf := oauth2.Config{
		ClientID:     APP_ID,
		ClientSecret: APP_KEY,
		RedirectURL:  REDIRECT_URL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://oauth.vk.com/authorize",
			TokenURL: "https://oauth.vk.com/access_token",
		},
		Scopes: []string{"email"},
	}

	if code == "" {
		return nil, 400, fmt.Errorf("wrong code")
	}

	token, err := conf.Exchange(ctx, code)
	fmt.Println("Token: ", token)
	if err != nil {
		return nil, 400, fmt.Errorf("cannot exchange")
	}

	fmt.Println("TOKEN OK")

	client := conf.Client(ctx, token)
	resp, err := client.Get(fmt.Sprintf(API_URL, token.AccessToken))
	if err != nil {
		return nil, 400, fmt.Errorf("cannot request data")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, fmt.Errorf("cannot read buffer")
	}

	data := &Response{}
	json.Unmarshal(body, data)

	birthdayTime, err := time.Parse("2006-01-02", data.Response[0].BirthDate)
	if err != nil {
		return nil, 500, fmt.Errorf("bad BirthDate")
	}

	vkUser := &api.VKUser{
		FirstName: data.Response[0].Name,
		Surname:   data.Response[0].LastName,
		Gender:    domain_models.GetGenderTypeInt(data.Response[0].Sex),
		Birthday:  birthdayTime,
		VKId:      uint32(data.Response[0].VKId),
	}

	fmt.Println("name: ", vkUser.FirstName, ", surname: ", vkUser.Surname, ", gender: ", vkUser.Gender, ", vkId: ", vkUser.VKId)
	return vkUser, 200, nil
}

// SignupVK handles user signup VK.
// @Summary User signup VK
// @Description Handles user signup VK.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Param invite_by query string false "invite_by value"
// @Success 200 {object} response.Response "Signup successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 500 {object} response.ErrorResponse "Failed to add user"
// @Router /api/v1/auth/auth-vk/signupVK [get]
func (ah *OAuthHandler) SignupVK(w http.ResponseWriter, r *http.Request) {
	newUser, status, err := AuthVK(w, r)
	if err != nil {
		response.HandleError(w, status, err.Error())
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

/*
// LoginVK handles user login.
// @Summary LoginVK User
// @Description LoginVK Handles user.
// @Tags auth-vk
// @Accept json
// @Produce json
// @Param code query string true "code from oauth"
// @Param invite_by query string false "invite_by param"
// @Success 200 {object} response.Response "Login successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Failed to create session"
// @Router /api/v1/auth/auth-vk/loginVK [get]
func (ah *OAuthHandler) LoginVK(w http.ResponseWriter, r *http.Request) {
		//userVK, status, err := AuthVK(w, r)
		//if err != nil {
		//	response.HandleError(w, status, err.Error())
		//	return
		//}
	userVK := &api.VKUser{
		FirstName: "sergey",
		Surname:   "fed",
		Gender:    domain_models.Male,
		VKId:      23435456,
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
*/
