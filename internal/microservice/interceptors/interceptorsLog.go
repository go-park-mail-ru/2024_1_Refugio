package interceptors

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

// InterceptorsLogger provides a structure for logging using Logrus.
type InterceptorsLogger struct {
	InterceptorsLogger *logrus.Logger
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
	output = strings.Replace(output, "%requestID%", entry.Data["requestID"].([]string)[0], 1)
	output = strings.Replace(output, "%time%", time.Now().Format("2006-01-02 15:04:05"), 1)
	output = strings.Replace(output, "%user-agent%", entry.Data["user-agent"].([]string)[0], 1)
	output = strings.Replace(output, "%FullMethod%", entry.Data["FullMethod"].(string), 1)
	output = strings.Replace(output, "%mode%", entry.Data["mode"].(string), 1)

	return []byte(output), nil
}

// InitializationAccessLogInterceptor initializes the logger for accessing resources.
func InitializationAccessLogInterceptor(f *os.File) *InterceptorsLogger {
	log := new(InterceptorsLogger)
	log.InterceptorsLogger = &logrus.Logger{
		Out:   io.MultiWriter(f, os.Stdout),
		Level: logrus.InfoLevel,
		Formatter: &Formatter{
			LogFormat:     "[%lvl%]: %time% - %msg% requestID=%requestID% work_time=%work_time% user-agent=%user-agent% FullMethod=%FullMethod% mode=%mode%\n",
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
