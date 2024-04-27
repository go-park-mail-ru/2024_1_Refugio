package http

import (
	"google.golang.org/grpc/metadata"
	auth_proto "mail/internal/microservice/auth/proto"
	"mail/internal/microservice/models/proto_converters"
	converters "mail/internal/models/delivery_converters"
	api "mail/internal/models/delivery_models"
	"mail/internal/models/microservice_ports"
	"mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"net/http"

	"mail/internal/pkg/utils/connect_microservice"
)

var (
	QHandler                        = &QuestionHandler{}
	requestIDContextKey interface{} = "requestid"
)

// QuestionHandler handles user-related HTTP requests.
type QuestionHandler struct {
	Sessions domainSession.SessionsManager
}

// Get all questions.
// @Summary Get questions
// @Description Get Handles questions.
// @Tags question
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "Login successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse "Failed to create session"
// @Router /api/v1/questions [get]
func (qh *QuestionHandler) GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.QuestionService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()

	authServiceClient := auth_proto.NewAuthServiceClient(conn)
	questionProto, errStatus := authServiceClient.Login(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&auth_proto.LoginRequest{},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusUnauthorized, "Login failed")
		return
	}

	questionsCore := proto_converters.EmailsConvertProtoInCore(questionProto)
	questionsApi := make([]*api.Question, 0, len(questionsCore))
	for _, question := range questionsCore {
		questionsApi = append(questionsApi, converters.QuestionConvertCoreInApi(*question))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"questions": questionsApi})
}
