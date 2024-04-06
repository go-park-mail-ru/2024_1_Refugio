package middleware

import (
	"github.com/sirupsen/logrus"
	"mail/pkg/delivery"
	"mail/pkg/delivery/session"
	"mail/pkg/domain/logger"
	"net/http"
	"time"
)

type Logger struct {
	Logger *logger.LogrusLogger
}

var Log = Logger{}
var requestIDContextKey interface{} = "requestid"

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// AuthMiddleware is a middleware to check user authentication using cookies.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, ok := r.Context().Value(requestIDContextKey).(string)
		if !ok {
			requestID = "none"
		}
		_, err := session.GlobalSeaaionManager.Check(r, requestID)
		if err != nil {
			delivery.HandleError(w, http.StatusUnauthorized, "Not Authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (log *Logger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		id, ok := r.Context().Value(requestIDContextKey).(string)
		if !ok {
			id = "none"
		}

		statusCode := lrw.statusCode
		en := log.Logger.LogrusLogger.WithFields(logrus.Fields{
			"method":     r.Method,
			"work_time":  time.Since(start),
			"URL":        r.URL.Path,
			"mode":       "[access_log]",
			"StatusCode": statusCode,
			"requestID":  id,
		})
		switch {
		case statusCode >= 200 && statusCode <= 207:
			en.Info("StatusOK")
		case statusCode >= 400 && statusCode <= 451:
			en.Warning("Client-side HTTP error")
		case statusCode >= 500 && statusCode <= 510:
			en.Error("StatusServerError")
		}

	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				delivery.HandleError(w, http.StatusInternalServerError, "StatusServerError")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
