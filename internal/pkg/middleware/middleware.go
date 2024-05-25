package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"mail/internal/monitoring"
	"mail/internal/pkg/logger"
	"mail/internal/pkg/session"

	response "mail/internal/models/response"
)

type Logger struct {
	Logger  *logger.LogrusLogger
	Metrics *monitoring.PrometheusMetrics
}

var (
	requestIDContextKey interface{} = "requestid"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// AuthMiddleware is a middleware to check user authentication using cookies.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := session.GlobalSessionManager.Check(r, r.Context())
		if err != nil {
			response.HandleError(w, http.StatusUnauthorized, "Not Authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (log *Logger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := NewLoggingResponseWriter(w)

		id, ok := r.Context().Value(requestIDContextKey).(string)
		if !ok {
			id = "none"
		}

		f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Failed to create logfile" + "log.txt")
		}
		defer f.Close()

		c := context.WithValue(r.Context(), "logger", logger.InitializationBdLog(f))
		ctx := context.WithValue(c, "requestID", id)

		req := r.WithContext(ctx)
		method := r.Method
		path := r.RequestURI
		re := regexp.MustCompile(`\d+`)
		customPath := strings.Split(re.ReplaceAllString(path, "1"), "?")[0]
		fmt.Println("customPath = ", customPath)

		log.Metrics.Hits.WithLabelValues(customPath, method).Inc()
		next.ServeHTTP(lrw, req)

		timing := time.Since(start)

		statusCode := lrw.statusCode
		en := log.Logger.LogrusLogger.WithFields(logrus.Fields{
			"method":     r.Method,
			"work_time":  timing,
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

		log.Metrics.Duration.WithLabelValues(strconv.Itoa(statusCode), customPath, method).Observe(timing.Seconds())
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				response.HandleError(w, http.StatusInternalServerError, "StatusServerError")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
