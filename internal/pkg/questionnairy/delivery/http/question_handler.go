package http

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"google.golang.org/grpc/metadata"
	domain "mail/internal/microservice/models/domain_models"
	"mail/internal/microservice/models/proto_converters"
	question_proto "mail/internal/microservice/questionnaire/proto"
	converters "mail/internal/models/delivery_converters"
	api "mail/internal/models/delivery_models"
	"mail/internal/models/microservice_ports"
	"mail/internal/models/response"
	domainSession "mail/internal/pkg/session/interface"
	"net/http"

	"mail/internal/pkg/utils/connect_microservice"
	"mail/internal/pkg/utils/sanitize"
	"mail/internal/pkg/utils/validators"
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
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Get questions successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to get questions"
// @Router /api/v1/questions [get]
func (qh *QuestionHandler) GetAllQuestions(w http.ResponseWriter, r *http.Request) {
	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.QuestionService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()

	questionServiceClient := question_proto.NewQuestionServiceClient(conn)
	questionProto, errStatus := questionServiceClient.GetQuestions(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&question_proto.GetQuestionsRequest{},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Get questions failed")
		return
	}

	questionsCore := make([]*domain.Question, 0, len(questionProto.Questions))
	for _, question := range questionProto.Questions {
		questionsCore = append(questionsCore, proto_converters.QuestionConvertProtoInCore(*question))
	}

	questionsApi := make([]*api.Question, 0, len(questionsCore))
	for _, question := range questionsCore {
		questionsApi = append(questionsApi, converters.QuestionConvertCoreInApi(*question))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"questions": questionsApi})
}

// Add question.
// @Summary Add question
// @Description Add Handles question.
// @Tags question
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param question body response.QuestionSwag true "Question message in JSON format"
// @Success 200 {object} response.Response "Add question successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to add question"
// @Router /api/v1/questions [post]
func (qh *QuestionHandler) AddQuestion(w http.ResponseWriter, r *http.Request) {
	var newQuestion api.Question
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newQuestion)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newQuestion.Text = sanitize.SanitizeString(newQuestion.Text)
	newQuestion.MinText = sanitize.SanitizeString(newQuestion.MinText)
	newQuestion.MaxText = sanitize.SanitizeString(newQuestion.MaxText)
	newQuestion.DopQuestion = sanitize.SanitizeString(newQuestion.DopQuestion)

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.QuestionService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()

	question := converters.QuestionConvertApiInCore(newQuestion)

	questionServiceClient := question_proto.NewQuestionServiceClient(conn)
	questionProto, errStatus := questionServiceClient.AddQuestion(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&question_proto.AddQuestionRequest{Question: proto_converters.QuestionConvertCoreInProto(*question)},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Add question failed")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": questionProto.Status})
}

// Add answer.
// @Summary Add answer
// @Description Add Handles answer.
// @Tags question
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Param question body response.AnswerSwag true "Answer message in JSON format"
// @Success 200 {object} response.Response "Add answer successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to add answer"
// @Router /api/v1/answers [post]
func (qh *QuestionHandler) AddAnswer(w http.ResponseWriter, r *http.Request) {
	var newAnswer api.Answer
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	err := json.NewDecoder(r.Body).Decode(&newAnswer)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newAnswer.Login = sanitize.SanitizeString(newAnswer.Login)
	newAnswer.Text = sanitize.SanitizeString(newAnswer.Text)

	if validators.IsValidEmailFormat(newAnswer.Login) {
		response.HandleError(w, http.StatusInternalServerError, "Login failed")
		return
	}

	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.QuestionService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()

	answer := converters.AnswerConvertApiInCore(newAnswer)

	questionServiceClient := question_proto.NewQuestionServiceClient(conn)
	answerProto, errStatus := questionServiceClient.AddAnswer(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&question_proto.AddAnswerRequest{Answer: proto_converters.AnswerConvertCoreInProto(*answer)},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Add answer failed")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": answerProto.Status})
}

// Get statistics.
// @Summary Get statistics
// @Description Get Handles statistics.
// @Tags question
// @Accept json
// @Produce json
// @Param X-Csrf-Token header string true "CSRF Token"
// @Success 200 {object} response.Response "Get statistics successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request body"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Failed to get statistics"
// @Router /api/v1/statistics [get]
func (qh *QuestionHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	conn, err := connect_microservice.OpenGRPCConnection(microservice_ports.GetPorts(microservice_ports.QuestionService))
	if err != nil {
		response.HandleError(w, http.StatusInternalServerError, "connection fail")
		return
	}
	defer conn.Close()

	questionServiceClient := question_proto.NewQuestionServiceClient(conn)
	statisticsProto, errStatus := questionServiceClient.GetStatistic(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&question_proto.GetStatisticRequest{},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Statistics failed")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"statistic": statisticsProto.Statistics})
}
