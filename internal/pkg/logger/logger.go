package logger

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// LogrusLogger provides a structure for logging using Logrus.
type LogrusLogger struct {
	LogrusLogger *logrus.Logger
}

// Formatter defines the formatting of logs.
type Formatter struct {
	LogFormat     string
	ForceColors   bool
	ColorInfo     *color.Color
	ColorWarning  *color.Color
	ColorError    *color.Color
	ColorCritical *color.Color
	ColorDefault  *color.Color
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
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

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Println("Error loc time")
	}
	time.Local = loc

	output = strings.Replace(output, "%msg%", entry.Message, 1)
	output = strings.Replace(output, "%work_time%", entry.Data["work_time"].(time.Duration).String(), 1)
	output = strings.Replace(output, "%mode%", entry.Data["mode"].(string), 1)
	output = strings.Replace(output, "%requestID%", entry.Data["requestID"].(string), 1)
	output = strings.Replace(output, "%time%", time.Now().Format("2006-01-02 15:04:05"), 1)

	if entry.Data["mode"].(string) == "[access_log]" {
		output = strings.Replace(output, "%port%", "8080", 1)
		output = strings.Replace(output, "%host%", "localhost", 1)
		output = strings.Replace(output, "%URL%", entry.Data["URL"].(string), 1)
		output = strings.Replace(output, "%method%", entry.Data["method"].(string), 1)
		output = strings.Replace(output, "%StatusCode%", strconv.Itoa(entry.Data["StatusCode"].(int)), 1)
	} else if entry.Data["mode"].(string) == "[db_log]" {
		output = strings.Replace(output, "%query%", entry.Data["query"].(string), 1)
		output = strings.Replace(output, "%args%", entry.Data["args"].(string), 1)
	}
	return []byte(output), nil
}

// InitializationBdLog initializes the logger to work with the database.
func InitializationBdLog(f *os.File) *LogrusLogger {
	log := new(LogrusLogger)
	log.LogrusLogger = &logrus.Logger{
		Out:   io.MultiWriter(f, os.Stdout),
		Level: logrus.InfoLevel,
		Formatter: &Formatter{
			LogFormat:     "[%lvl%]: %time% - %msg% requestID=%requestID% work_time=%work_time% mode=%mode% query=%query% arguments=%args%\n",
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

// InitializationAccessLog initializes the logger for accessing resources.
func InitializationAccessLog(f *os.File) *LogrusLogger {
	log := new(LogrusLogger)
	log.LogrusLogger = &logrus.Logger{
		Out:   io.MultiWriter(f, os.Stdout),
		Level: logrus.InfoLevel,
		Formatter: &Formatter{
			LogFormat:     "[%lvl%]: %time% - %msg% requestID=%requestID% method=%method% StatusCode=%StatusCode% host=%host% port=%port% URL=%URL% work_time=%work_time% mode=%mode%\n",
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

// InitializationEmptyLog initializes an empty logger Logrus.
func InitializationEmptyLog() *LogrusLogger {
	log := new(LogrusLogger)
	log.LogrusLogger = &logrus.Logger{}
	return log
}

// DbLog logs information about database requests.
func (log *LogrusLogger) DbLog(query, requestID string, start time.Time, err *error, args []interface{}) {
	requestID = GetRequestIDString(requestID)
	resArgs := "{ "
	for _, a := range args {
		resArgs += fmt.Sprintf("%#v, ", a)
	}
	resArgs += "}"
	en := log.LogrusLogger.WithFields(logrus.Fields{
		"work_time": time.Since(start),
		"requestID": requestID,
		"mode":      "[db_log]",
		"query":     query,
		"args":      resArgs,
	})

	if *err != nil {
		en.Error("StatusServerError")
	} else {
		en.Info("StatusOK")
	}
}

// GetRequestIDString converts the request ID Value to a string.
func GetRequestIDString(requestIDValue interface{}) string {
	if requestIDValue != nil {
		requestIDString, ok := requestIDValue.(string)
		if !ok {
			requestIDString = fmt.Sprintf("%v", requestIDValue)
			requestIDString = strings.Trim(requestIDString, "[]")
		}

		return requestIDString
	}

	return ""
}
