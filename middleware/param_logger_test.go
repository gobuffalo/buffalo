package middleware

import (
	"net/url"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/willie"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_maskSecrets(t *testing.T) {
	r := require.New(t)
	pl := parameterLogger{}

	filteredForm := pl.maskSecrets(url.Values{
		"FirstName":            []string{"Antonio"},
		"MiddleName":           []string{"José"},
		"LastName":             []string{"Pagano"},
		"Password":             []string{"Secret!"},
		"password":             []string{"Other"},
		"pAssWorD":             []string{"Weird one"},
		"PasswordConfirmation": []string{"Secret!"},

		"SomeCVC": []string{"Untouched"},
	})

	r.Equal(filteredForm.Get("Password"), filteredIndicator[0])
	r.Equal(filteredForm.Get("password"), filteredIndicator[0])
	r.Equal(filteredForm.Get("pAssWorD"), filteredIndicator[0])
	r.Equal(filteredForm.Get("PasswordConfirmation"), filteredIndicator[0])
	r.Equal(filteredForm.Get("LastName"), "Pagano")
	r.Equal(filteredForm.Get("SomeCVC"), "Untouched")
}

func Test_maskSecretsCustom(t *testing.T) {
	r := require.New(t)
	pl := parameterLogger{
		blacklist: []string{
			"FirstName", "LastName", "MiddleName",
		},
	}

	filteredForm := pl.maskSecrets(url.Values{
		"FirstName":            []string{"Antonio"},
		"MiddleName":           []string{"José"},
		"LastName":             []string{"Pagano"},
		"Password":             []string{"Secret!"},
		"password":             []string{"Other"},
		"pAssWorD":             []string{"Weird one"},
		"PasswordConfirmation": []string{"Secret!"},

		"SomeCVC": []string{"Untouched"},
	})

	r.Equal(filteredForm.Get("Password"), "Secret!")
	r.Equal(filteredForm.Get("password"), "Other")
	r.Equal(filteredForm.Get("LastName"), filteredIndicator[0])
	r.Equal(filteredForm.Get("SomeCVC"), "Untouched")
}

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
