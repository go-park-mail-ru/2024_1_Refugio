package logger

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"
	_ "time/tzdata"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestInitializationBdLog(t *testing.T) {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile in session_repo" + "log.txt")
	}
	defer f.Close()
	expectedLog := new(LogrusLogger)
	expectedLog.LogrusLogger = &logrus.Logger{
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

	log := InitializationBdLog(f)

	assert.Equal(t, expectedLog, log)
}

func TestInitializationAccesLog(t *testing.T) {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile in session_repo" + "log.txt")
	}
	defer f.Close()
	expectedLog := new(LogrusLogger)
	expectedLog.LogrusLogger = &logrus.Logger{
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

	log := InitializationAccessLog(f)

	assert.Equal(t, *expectedLog, *log)
}

func TestDbLog(t *testing.T) {
	testlogger, hook := test.NewNullLogger()
	f := &Formatter{
		LogFormat:     "[%lvl%]: %time% - %msg% requestID=%requestID% work_time=%work_time% mode=%mode% query=%query% arguments=%args%\n",
		ForceColors:   true,
		ColorInfo:     color.New(color.FgBlue),
		ColorWarning:  color.New(color.FgYellow),
		ColorError:    color.New(color.FgRed),
		ColorCritical: color.New(color.BgRed, color.FgWhite),
		ColorDefault:  color.New(color.FgWhite),
	}
	testlogger.SetFormatter(f)
	testlogger.SetLevel(logrus.InfoLevel)

	logger := new(LogrusLogger)
	logger.LogrusLogger = testlogger

	query := "TestQuery"
	requestID := "test_request"
	start := time.Now()
	args := []interface{}{}

	t.Run("DbLogSuccessfully", func(t *testing.T) {
		var err error
		logger.DbLog(query, requestID, start, &err, args)

		expectedData := logrus.Fields{
			"work_time": hook.LastEntry().Data["work_time"],
			"requestID": requestID,
			"mode":      "[db_log]",
			"query":     query,
			"args":      "{ }",
		}

		assert.Equal(t, "StatusOK", hook.LastEntry().Message)
		assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
		assert.Equal(t, expectedData, hook.LastEntry().Data)
	})

	t.Run("DbLogFail", func(t *testing.T) {
		err := fmt.Errorf("Error query=%v", query)

		logger.DbLog(query, requestID, start, &err, args)

		expectedData := logrus.Fields{
			"work_time": hook.LastEntry().Data["work_time"],
			"requestID": requestID,
			"mode":      "[db_log]",
			"query":     query,
			"args":      "{ }",
		}

		assert.Equal(t, "StatusServerError", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		assert.Equal(t, expectedData, hook.LastEntry().Data)
	})
}

func TestFormat(t *testing.T) {
	requestID := "test_request"
	start := time.Now()

	t.Run("FormatDbLog", func(t *testing.T) {
		f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Failed to create logfile in session_repo" + "log.txt")
		}
		defer f.Close()
		log := InitializationBdLog(f)

		en := log.LogrusLogger.WithFields(logrus.Fields{
			"work_time": time.Since(start),
			"requestID": requestID,
			"mode":      "[db_log]",
			"query":     "TestQuery",
			"args":      "{ }",
		})
		en.Level = logrus.InfoLevel

		b, err := log.LogrusLogger.Formatter.Format(en)
		if err != nil {
			return
		}

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Println("Error loc time")
		}
		time.Local = loc

		expectedLog := fmt.Sprintf("[INFO]: %v -  requestID=test_request work_time=%v mode=[db_log] query=TestQuery arguments={ }\n", time.Now().Format("2006-01-02 15:04:05"), en.Data["work_time"].(time.Duration).String())
		assert.NoError(t, err)
		assert.Equal(t, expectedLog, string(b))
	})

	t.Run("FormatAccessLog", func(t *testing.T) {
		f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Failed to create logfile in session_repo" + "log.txt")
		}
		defer f.Close()
		log := InitializationAccessLog(f)

		en := log.LogrusLogger.WithFields(logrus.Fields{
			"method":     "GET",
			"work_time":  time.Since(start),
			"URL":        "/",
			"mode":       "[access_log]",
			"StatusCode": 200,
			"requestID":  requestID,
		})
		en.Level = logrus.InfoLevel

		b, err := log.LogrusLogger.Formatter.Format(en)
		if err != nil {
			return
		}

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Println("Error loc time")
		}
		time.Local = loc

		expectedLog := fmt.Sprintf("[INFO]: %v -  requestID=test_request method=GET StatusCode=200 host=localhost port=8080 URL=/ work_time=%v mode=[access_log]\n", time.Now().Format("2006-01-02 15:04:05"), en.Data["work_time"].(time.Duration).String())
		assert.NoError(t, err)
		assert.Equal(t, expectedLog, string(b))
	})
}
