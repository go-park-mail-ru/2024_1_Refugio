package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"mail/pkg/delivery"
	"mail/pkg/delivery/session"
	"net/http"
	"time"
)

type LogrusLogger struct {
	LogrusLogger *logrus.Entry
}

func InitializationAcceslog(port int) *LogrusLogger {
	Logrus := new(LogrusLogger)
	Logrus.LogrusLogger = logrus.WithFields(logrus.Fields{
		"logger": "Logrus",
		"host":   "localhost",
		"port":   port,
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	return Logrus
}

// AuthMiddleware is a middleware to check user authentication using cookies.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AuthMiddleware", r.URL.Path)
		_, err := session.GlobalSeaaionManager.Check(r)
		if err != nil {
			delivery.HandleError(w, http.StatusUnauthorized, "Not Authorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (ac *LogrusLogger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		ac.LogrusLogger.WithFields(logrus.Fields{
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"work_time":   time.Since(start),
			"URL":         r.URL.Path,
			"mode":        "[access_log]",
		}).Info("AccessLogMiddleware")
	})
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AccessLogMiddleware", r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Printf("[%s] %s, %s %s\n", r.Method, r.RemoteAddr, r.URL.Path, time.Since(start))
	})
}

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("PanicMiddleware", r.URL.Path)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
