package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"mail/pkg/delivery"
	"mail/pkg/delivery/session"
	"net/http"
	"os"
	"strings"
	"time"
)

type LogrusLogger struct {
	LogrusLogger *logrus.Logger
}

type Formatter struct {
	LogFormat string
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	output = strings.Replace(output, "%time%", entry.Time.Format("2006-01-02 15:04:05"), 1)
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	output = strings.Replace(output, "%lvl%", strings.ToUpper(entry.Level.String()), 1)
	output = strings.Replace(output, "%logger%", "Logrus", 1)
	output = strings.Replace(output, "%host%", "localhost", 1)
	output = strings.Replace(output, "%port%", "8080", 1)
	output = strings.Replace(output, "%URL%", entry.Data["URL"].(string), 1)
	output = strings.Replace(output, "%method%", entry.Data["method"].(string), 1)
	output = strings.Replace(output, "%work_time%", entry.Data["work_time"].(time.Duration).String(), 1)
	//output = strings.Replace(output, "%remote_addr%", entry.Data["remote_addr"].(string), 1)
	output = strings.Replace(output, "%access_log%", entry.Data["mode"].(string), 1)
	return []byte(output), nil
}

func InitializationAcceslog() *LogrusLogger {
	f := &Formatter{LogFormat: "[%lvl%]: %time% - %msg% (%logger%) host=%host% port=%port% URL=%URL% method=%method% work_time=%work_time% remote_addr=%remote_addr% access_log=%access_log%\n"}
	log := new(LogrusLogger)
	log.LogrusLogger = &logrus.Logger{}
	log.LogrusLogger.SetFormatter(f)
	log.LogrusLogger.SetLevel(logrus.InfoLevel)
	log.LogrusLogger.Out = os.Stdout
	return log
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

func (log *LogrusLogger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		en := log.LogrusLogger.WithFields(logrus.Fields{
			"method":    r.Method,
			"work_time": time.Since(start),
			"URL":       r.URL.Path,
			"mode":      "[access_log]",
		})
		/*en := logrus.WithFields(logrus.Fields{
		  "method":    r.Method,
		  "work_time": time.Since(start),
		  "URL":       r.URL.Path,
		  "mode":      "[access_log]",
		})*/
		en.Info("AccessLogMiddleware")
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
