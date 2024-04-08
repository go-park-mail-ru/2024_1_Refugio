package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
	"time"
)

func TestInitializationBdLog(t *testing.T) {
	expectedLog := new(LogrusLogger)
	expectedLog.LogrusLogger = &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout),
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

	log := InitializationBdLog()

	assert.Equal(t, expectedLog, log)
}

func TestInitializationAccesLog(t *testing.T) {
	expectedLog := new(LogrusLogger)
	expectedLog.LogrusLogger = &logrus.Logger{
		Out:   io.MultiWriter(os.Stdout),
		Level: logrus.InfoLevel,
		Formatter: &Formatter{
			LogFormat:     "[%lvl%]: %time% - %msg% method=%method% StatusCode=%StatusCode% requestID=%requestID% host=%host% port=%port% URL=%URL% work_time=%work_time% mode=%mode%\n",
			ForceColors:   true,
			ColorInfo:     color.New(color.FgBlue),
			ColorWarning:  color.New(color.FgYellow),
			ColorError:    color.New(color.FgRed),
			ColorCritical: color.New(color.BgRed, color.FgWhite),
			ColorDefault:  color.New(color.FgWhite),
		},
	}

	log := InitializationAccesLog()

	assert.Equal(t, expectedLog, log)
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
		log := InitializationBdLog()

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

		expectedLog := fmt.Sprintf("[INFO]: 0001-01-01 00:00:00 -  requestID=test_request work_time=%v mode=[db_log] query=TestQuery arguments={ }\n", en.Data["work_time"].(time.Duration).String())
		assert.NoError(t, err)
		assert.Equal(t, expectedLog, string(b))
	})

	t.Run("FormatAccessLog", func(t *testing.T) {
		log := InitializationAccesLog()

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

		fmt.Println(string(b))

		expectedLog := fmt.Sprintf("[INFO]: 0001-01-01 00:00:00 -  method=GET StatusCode=200 requestID=test_request host=localhost port=8080 URL=/ work_time=%v mode=[access_log]\n", en.Data["work_time"].(time.Duration).String())
		assert.NoError(t, err)
		assert.Equal(t, expectedLog, string(b))
	})
}
