package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicMiddleware(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	t.Run("WithoutCookie", func(t *testing.T) {
		reqWithoutCookie := httptest.NewRequest("GET", "/api/v1/auth/verify-auth", nil)
		recWithoutCookie := httptest.NewRecorder()
		PanicMiddleware(fakeHandler).ServeHTTP(recWithoutCookie, reqWithoutCookie)
		assert.Equal(t, http.StatusOK, recWithoutCookie.Code)
	})

}

func TestWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	lrw := NewLoggingResponseWriter(w)
	lrw.WriteHeader(200)

	expectedLrw := new(loggingResponseWriter)
	expectedLrw.statusCode = 200
	expectedLrw.ResponseWriter = w

	assert.Equal(t, expectedLrw, lrw)
}

func TestAuthMiddleware(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("WithoutCookie", func(t *testing.T) {
		reqWithoutCookie := httptest.NewRequest("GET", "/api/v1/auth/verify-auth", nil)
		recWithoutCookie := httptest.NewRecorder()
		AuthMiddleware(fakeHandler).ServeHTTP(recWithoutCookie, reqWithoutCookie)
		assert.Equal(t, http.StatusUnauthorized, recWithoutCookie.Code)
	})

	t.Run("WithInvalidCookie", func(t *testing.T) {
		reqWithInvalidCookie := httptest.NewRequest("GET", "/api/v1/auth/verify-auth", nil)
		reqWithInvalidCookie.AddCookie(&http.Cookie{Name: "session_id", Value: ""})
		recWithInvalidCookie := httptest.NewRecorder()
		AuthMiddleware(fakeHandler).ServeHTTP(recWithInvalidCookie, reqWithInvalidCookie)
		assert.Equal(t, http.StatusUnauthorized, recWithInvalidCookie.Code)
	})

	/*
		t.Run("StatusOK", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockSessionServiceClient := mock.NewMockSessionServiceClient(ctrl)
			sessionManager := session.NewSessionsManager(mockSessionServiceClient)
			session.InitializationGlobalSessionManager(sessionManager)
			expectedSession := session_proto.Session{
				SessionId: "session_id",
				UserId:    1,
				Device:    "desktop",
				LifeTime:  3600,
				CsrfToken: "csrf_token",
			}

			requestID := "testID"

			f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Println("Failed to create logfile" + "log.txt")
			}
			defer f.Close()

			c := context.WithValue(context.Background(), "logger", logger.InitializationBdLog(f))
			ctx := context.WithValue(c, "requestID", requestID)

			req, _ := http.NewRequest("GET", "/api/v1/auth/verify-auth", nil)
			req.WithContext(ctx)
			req.Header.Set("X-CSRF-Token", "csrf_token")
			req.AddCookie(&http.Cookie{Name: "session_id", Value: "session_id"})

			mockSessionServiceClient.EXPECT().
				GetSession(gomock.Any(), gomock.Any()).
				Return(&session_proto.GetSessionReply{Session: &expectedSession}, nil).
				Times(1)

			recW := httptest.NewRecorder()
			AuthMiddleware(fakeHandler).ServeHTTP(recW, req)
			assert.Equal(t, http.StatusOK, recW.Code)
		})
	*/
}

/*
func TestAccessLogMiddleware(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	testlogger, hook := test.NewNullLogger()

	f := &logger.Formatter{
		LogFormat:     "[%lvl%]: %time% - %msg% method=%method% StatusCode=%StatusCode% requestID=%requestID% host=%host% port=%port% URL=%URL% work_time=%work_time% remote_addr=%remote_addr% access_log=%access_log%\n",
		ForceColors:   true,
		ColorInfo:     color.New(color.FgBlue),
		ColorWarning:  color.New(color.FgYellow),
		ColorError:    color.New(color.FgRed),
		ColorCritical: color.New(color.BgRed, color.FgWhite),
		ColorDefault:  color.New(color.FgWhite),
	}

	testlogger.SetFormatter(f)
	testlogger.SetLevel(logrus.InfoLevel)
	log := new(logger.LogrusLogger)
	log.LogrusLogger = testlogger
	Logrus := new(Logger)
	Logrus.Logger = log

	Logrus.AccessLogMiddleware(fakeHandler).ServeHTTP(rec, req)

	if hook.LastEntry().Message != "StatusOK" {
		t.Errorf("Bad input handled incorrectly")
	}
	expectedData := logrus.Fields{
		"method":     "GET",
		"work_time":  hook.LastEntry().Data["work_time"],
		"URL":        "/",
		"mode":       "[access_log]",
		"StatusCode": 200,
		"requestID":  hook.LastEntry().Data["requestID"].(string),
	}

	assert.Equal(t, "StatusOK", hook.LastEntry().Message)
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, expectedData, hook.LastEntry().Data)
}
*/
