package middleware

import (
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"io"
	"mail/pkg/delivery"
	"mail/pkg/delivery/session"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type LogrusLogger struct {
	LogrusLogger *logrus.Logger
}

type Formatter struct {
	LogFormat     string
	ForceColors   bool
	ColorInfo     *color.Color
	ColorWarning  *color.Color
	ColorError    *color.Color
	ColorCritical *color.Color
	ColorDefault  *color.Color
}

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

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	//var levelColor int
	output := f.LogFormat
	if f.ForceColors {
		switch entry.Level {
		case logrus.InfoLevel:
			output = strings.Replace(output, "%lvl%", f.ColorInfo.Sprintf("%s", strings.ToUpper(entry.Level.String())), 1) //blue
		case logrus.WarnLevel:
			output = strings.Replace(output, "%lvl%", f.ColorWarning.Sprintf("%s", strings.ToUpper(entry.Level.String())), 1) //yellow
		case logrus.ErrorLevel:
			output = strings.Replace(output, "%lvl%", f.ColorError.Sprintf("%s", strings.ToUpper(entry.Level.String())), 1) //red
		case logrus.FatalLevel, logrus.PanicLevel:
			output = strings.Replace(output, "%lvl%", f.ColorCritical.Sprintf("%s", strings.ToUpper(entry.Level.String())), 1) //red background and white text
		default:
			output = strings.Replace(output, "%lvl%", f.ColorDefault.Sprintf("%s", strings.ToUpper(entry.Level.String())), 1) //white
		}
	} else {
		output = strings.Replace(output, "%lvl%", f.ColorDefault.Sprintf("%s", strings.ToUpper(entry.Level.String())), 1)
	}
	output = strings.Replace(output, "%msg%", entry.Message, 1)
	output = strings.Replace(output, "%port%", "8080", 1)
	output = strings.Replace(output, "%host%", "localhost", 1)
	output = strings.Replace(output, "%URL%", entry.Data["URL"].(string), 1)
	output = strings.Replace(output, "%method%", entry.Data["method"].(string), 1)
	output = strings.Replace(output, "%access_log%", entry.Data["mode"].(string), 1)
	output = strings.Replace(output, "%requestID%", entry.Data["requestID"].(string), 1)
	output = strings.Replace(output, "%time%", entry.Time.Format("2006-01-02 15:04:05"), 1)
	output = strings.Replace(output, "%StatusCode%", strconv.Itoa(entry.Data["StatusCode"].(int)), 1)
	output = strings.Replace(output, "%work_time%", entry.Data["work_time"].(time.Duration).String(), 1)
	return []byte(output), nil
}

func InitializationAcceslog(f *os.File) *LogrusLogger {
	log := new(LogrusLogger)
	log.LogrusLogger = &logrus.Logger{
		Out:   io.MultiWriter(f, os.Stdout), //os.Stdout,
		Level: logrus.InfoLevel,
		Formatter: &Formatter{
			LogFormat:     "[%lvl%]: %time% - %msg% method=%method% StatusCode=%StatusCode% requestID=%requestID% host=%host% port=%port% URL=%URL% work_time=%work_time% remote_addr=%remote_addr% access_log=%access_log%\n",
			ForceColors:   true,
			ColorInfo:     color.New(color.FgBlue),
			ColorWarning:  color.New(color.FgYellow),
			ColorError:    color.New(color.FgRed),
			ColorCritical: color.New(color.BgRed, color.FgWhite),
			ColorDefault:  color.New(color.FgWhite),
		},
	}

	return log
}

// AuthMiddleware is a middleware to check user authentication using cookies.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		id, ok := r.Context().Value(requestIDContextKey).(string)
		if !ok {
			id = "none"
		}

		statusCode := lrw.statusCode
		en := log.LogrusLogger.WithFields(logrus.Fields{
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
