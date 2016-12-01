package buffalo

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
	"github.com/markbates/going/defaults"
)

type Options struct {
	Env            string
	LogLevel       string
	Logger         Logger
	LogDir         string
	NotFound       http.Handler
	MethodOverride func(r *http.Request)
	// Store is the `github.com/gorilla/sessions` store used to back
	// the session. It defaults to use a cookie store and the ENV variable
	// `SESSION_SECRET`.
	SessionStore sessions.Store
	// SessionName is the name of the session cookie that is set. This defaults
	// to "_buffalo_session".
	SessionName string
	prefix      string
}

func NewOptions() Options {
	return optionsWithDefaults(Options{})
}

func optionsWithDefaults(opts Options) Options {
	opts.Env = defaults.String(opts.Env, defaults.String(os.Getenv("BUFFALO_ENV"), defaults.String(os.Getenv("GO_ENV"), "development")))
	opts.LogLevel = defaults.String(opts.LogLevel, "debug")
	pwd, _ := os.Getwd()
	opts.LogDir = defaults.String(opts.LogDir, filepath.Join(pwd, "logs"))
	if opts.SessionStore == nil {
		opts.SessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	}
	opts.SessionName = defaults.String(opts.SessionName, "_buffalo_session")
	return opts
}
