package http

import (
	"google.golang.org/grpc/metadata"
	"io"
	"net/http"

	"mail/internal/microservice/models/proto_converters"
	"mail/internal/models/response"
	"mail/internal/pkg/utils/sanitize"

	domain "mail/internal/microservice/models/domain_models"
	question_proto "mail/internal/microservice/questionnaire/proto"
	converters "mail/internal/models/delivery_converters"
	api "mail/internal/models/delivery_models"
	domainSession "mail/internal/pkg/session/interface"
)

var (
	QHandler = &QuestionHandler{}
)

// QuestionHandler handles user-related HTTP requests.
type QuestionHandler struct {
	Sessions              domainSession.SessionsManager
	QuestionServiceClient question_proto.QuestionServiceClient
}

// GetAllQuestions all questions.
// @Summary GetAllQuestions questions
// @Description GetAllQuestions Handles questions.
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
	questionProto, errStatus := qh.QuestionServiceClient.GetQuestions(
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
		questionsCore = append(questionsCore, proto_converters.QuestionConvertProtoInCore(question))
	}

	questionsApi := make([]*api.Question, 0, len(questionsCore))
	for _, question := range questionsCore {
		questionsApi = append(questionsApi, converters.QuestionConvertCoreInApi(*question))
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"questions": questionsApi})
}

// AddQuestion question.
// @Summary AddQuestion question
// @Description AddQuestion Handles question.
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newQuestion api.Question
	if err := newQuestion.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newQuestion.Text = sanitize.SanitizeString(newQuestion.Text)
	newQuestion.MinText = sanitize.SanitizeString(newQuestion.MinText)
	newQuestion.MaxText = sanitize.SanitizeString(newQuestion.MaxText)
	newQuestion.DopQuestion = sanitize.SanitizeString(newQuestion.DopQuestion)

	question := converters.QuestionConvertApiInCore(newQuestion)

	questionProto, errStatus := qh.QuestionServiceClient.AddQuestion(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&question_proto.AddQuestionRequest{Question: proto_converters.QuestionConvertCoreInProto(question)},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Add question failed")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": questionProto.Status})
}

// AddAnswer answer.
// @Summary AddAnswer answer
// @Description AddAnswer Handles answer.
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.HandleError(w, http.StatusBadRequest, "Invalid input body")
		return
	}
	var newAnswer api.Answer
	if err := newAnswer.UnmarshalJSON(body); err != nil {
		response.HandleError(w, http.StatusBadRequest, "Bad JSON in request")
		return
	}

	newAnswer.Text = sanitize.SanitizeString(newAnswer.Text)

	login, errLogin := qh.Sessions.GetLoginBySession(r, r.Context())
	if errLogin != nil {
		response.HandleError(w, http.StatusInternalServerError, "Login fail")
		return
	}
	newAnswer.Login = login

	answer := converters.AnswerConvertApiInCore(newAnswer)

	answerProto, errStatus := qh.QuestionServiceClient.AddAnswer(
		metadata.NewOutgoingContext(r.Context(),
			metadata.New(map[string]string{"requestID": r.Context().Value("requestID").(string)})),
		&question_proto.AddAnswerRequest{Answer: proto_converters.AnswerConvertCoreInProto(answer)},
	)
	if errStatus != nil {
		response.HandleError(w, http.StatusInternalServerError, "Add answer failed")
		return
	}

	response.HandleSuccess(w, http.StatusOK, map[string]interface{}{"Success": answerProto.Status})
}

// GetStatistics statistics.
// @Summary GetStatistics statistics
// @Description GetStatistics Handles statistics.
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
	statisticsProto, errStatus := qh.QuestionServiceClient.GetStatistic(
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
