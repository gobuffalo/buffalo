package buffalo

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/markbates/going/defaults"
)

// Options are used to configure and define how your application should run.
type Options struct {
	// Env is the "environment" in which the App is running. Default is "development".
	Env string
	// LogLevel defaults to "debug".
	LogLevel string
	// Logger to be used with the application. A default one is provided.
	Logger Logger
	// LogDir is the path to the directory for storing the JSON log files from the
	// default Logger
	LogDir         string
	NotFound       http.Handler
	MethodOverride http.HandlerFunc
	// SessionStore is the `github.com/gorilla/sessions` store used to back
	// the session. It defaults to use a cookie store and the ENV variable
	// `SESSION_SECRET`.
	SessionStore sessions.Store
	// SessionName is the name of the session cookie that is set. This defaults
	// to "_buffalo_session".
	SessionName string
	prefix      string
}

// NewOptions returns a new Options instance with sensible defaults
func NewOptions() Options {
	return optionsWithDefaults(Options{})
}

func optionsWithDefaults(opts Options) Options {
	opts.Env = defaults.String(opts.Env, defaults.String(os.Getenv("GO_ENV"), "development"))
	opts.LogLevel = defaults.String(opts.LogLevel, "debug")
	pwd, _ := os.Getwd()
	opts.LogDir = defaults.String(opts.LogDir, filepath.Join(pwd, "logs"))
	if opts.SessionStore == nil {
		opts.SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	}
	opts.SessionName = defaults.String(opts.SessionName, "_buffalo_session")
	return opts
}
