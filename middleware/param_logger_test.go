package middleware

import (
	"net/url"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/willie"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var lastEntry *logrus.Entry

type testHook struct{}

func (th testHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (th testHook) Fire(entry *logrus.Entry) error {
	lastEntry = entry
	return nil
}

type testLogger struct {
	logrus.FieldLogger
}

func (l testLogger) WithField(s string, i interface{}) buffalo.Logger {
	return testLogger{l.FieldLogger.WithField(s, i)}
}

func (l testLogger) WithFields(m map[string]interface{}) buffalo.Logger {
	return testLogger{l.FieldLogger.WithFields(m)}
}

func newTestLogger() testLogger {
	l := logrus.New()
	l.AddHook(testHook{})
	l.Level, _ = logrus.ParseLevel("debug")

	return testLogger{l}
}

func Test_Logger(t *testing.T) {
	r := require.New(t)
	app := buffalo.New(buffalo.Options{})
	app.Use(ParameterLogger)
	app.Logger = newTestLogger()
	emptyHandler := func(c buffalo.Context) error {
		return nil
	}

	app.GET("/", emptyHandler)
	app.POST("/", emptyHandler)

	wi := willie.New(app)
	wi.HTML("/?param=value").Get()

	r.Contains(lastEntry.Data["params"], "{\"param\":[\"value\"]}")

	wi.HTML("/").Post(url.Values{
		"Password": []string{"123"},
		"Name":     []string{"Antonio"},
		"CVC":      []string{"123"},
	})

	r.Contains(lastEntry.Data["form"], "\"CVC\":[\"[FILTERED]\"]")
	r.Contains(lastEntry.Data["form"], "\"Name\":[\"Antonio\"]")
	r.Contains(lastEntry.Data["form"], "\"Password\":[\"[FILTERED]\"]")
}
