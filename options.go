package buffalo

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/logging"
	"github.com/gobuffalo/x/defaults"
	"github.com/gorilla/sessions"
	"github.com/markbates/oncer"
)

// Options are used to configure and define how your application should run.
type Options struct {
	Name string `json:"name"`
	// Addr is the bind address provided to http.Server. Default is "127.0.0.1:3000"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string `json:"addr"`
	// Host that this application will be available at. Default is "http://127.0.0.1:[$PORT|3000]".
	Host string `json:"host"`

	// Env is the "environment" in which the App is running. Default is "development".
	Env string `json:"env"`

	// LogLevel defaults to "debug". Deprecated use LogLvl instead
	LogLevel string `json:"log_level"`
	// LogLevl defaults to logger.DebugLvl.
	LogLvl logger.Level `json:"log_lvl"`
	// Logger to be used with the application. A default one is provided.
	Logger Logger `json:"-"`

	// MethodOverride allows for changing of the request method type. See the default
	// implementation at buffalo.MethodOverride
	MethodOverride http.HandlerFunc `json:"-"`

	// SessionStore is the `github.com/gorilla/sessions` store used to back
	// the session. It defaults to use a cookie store and the ENV variable
	// `SESSION_SECRET`.
	SessionStore sessions.Store `json:"-"`
	// SessionName is the name of the session cookie that is set. This defaults
	// to "_buffalo_session".
	SessionName string `json:"session_name"`

	// Worker implements the Worker interface and can process tasks in the background.
	// Default is "github.com/gobuffalo/worker.Simple.
	Worker worker.Worker `json:"-"`
	// WorkerOff tells App.Start() whether to start the Worker process or not. Default is "false".
	WorkerOff bool `json:"worker_off"`

	// PreHandlers are http.Handlers that are called between the http.Server
	// and the buffalo Application.
	PreHandlers []http.Handler `json:"-"`
	// PreWare takes an http.Handler and returns and http.Handler
	// and acts as a pseudo-middleware between the http.Server and
	// a Buffalo application.
	PreWares []PreWare `json:"-"`

	Prefix  string          `json:"prefix"`
	Context context.Context `json:"-"`

	cancel context.CancelFunc
}

// PreWare takes an http.Handler and returns and http.Handler
// and acts as a pseudo-middleware between the http.Server and
// a Buffalo application.
type PreWare func(http.Handler) http.Handler

// NewOptions returns a new Options instance with sensible defaults
func NewOptions() Options {
	return optionsWithDefaults(Options{})
}

func optionsWithDefaults(opts Options) Options {
	opts.Env = defaults.String(opts.Env, envy.Get("GO_ENV", "development"))
	opts.Name = defaults.String(opts.Name, "/")
	addr := "0.0.0.0"
	if opts.Env == "development" {
		addr = "127.0.0.1"
	}
	envAddr := envy.Get("ADDR", addr)

	if strings.HasPrefix(envAddr, "unix:") {
		// UNIX domain socket doesn't have a port
		opts.Addr = envAddr
	} else {
		// TCP case
		opts.Addr = defaults.String(opts.Addr, fmt.Sprintf("%s:%s", envAddr, envy.Get("PORT", "3000")))
	}

	if opts.PreWares == nil {
		opts.PreWares = []PreWare{}
	}
	if opts.PreHandlers == nil {
		opts.PreHandlers = []http.Handler{}
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}
	opts.Context, opts.cancel = context.WithCancel(opts.Context)

	if opts.Logger == nil {
		if lvl, err := envy.MustGet("LOG_LEVEL"); err == nil {
			opts.LogLvl, err = logger.ParseLevel(lvl)
			if err != nil {
				opts.LogLvl = logger.DebugLevel
			}
		}

		if len(opts.LogLevel) > 0 {
			var err error
			oncer.Deprecate(0, "github.com/gobuffalo/buffalo#Options.LogLevel", "Use github.com/gobuffalo/buffalo#Options.LogLvl instead.")
			opts.LogLvl, err = logger.ParseLevel(opts.LogLevel)
			if err != nil {
				opts.LogLvl = logger.DebugLevel
			}
		}
		if opts.LogLvl == 0 {
			opts.LogLvl = logger.DebugLevel
		}

		opts.Logger = logger.New(opts.LogLvl)
	}

	pop.SetLogger(func(level logging.Level, s string, args ...interface{}) {
		if !pop.Debug {
			return
		}

		l := opts.Logger
		if len(args) > 0 {
			for i, a := range args {
				l = l.WithField(fmt.Sprintf("$%d", i+1), a)
			}
		}

		if pop.Color {
			s = color.YellowString(s)
		}

		l.Debug(s)
	})

	if opts.SessionStore == nil {
		secret := envy.Get("SESSION_SECRET", "")
		// In production a SESSION_SECRET must be set!
		if secret == "" {
			if opts.Env == "development" || opts.Env == "test" {
				secret = "buffalo-secret"
			} else {
				opts.Logger.Warn("Unless you set SESSION_SECRET env variable, your session storage is not protected!")
			}
		}
		opts.SessionStore = sessions.NewCookieStore([]byte(secret))
	}
	if opts.Worker == nil {
		w := worker.NewSimpleWithContext(opts.Context)
		w.Logger = opts.Logger
		opts.Worker = w
	}
	opts.SessionName = defaults.String(opts.SessionName, "_buffalo_session")
	opts.Host = defaults.String(opts.Host, envy.Get("HOST", fmt.Sprintf("http://127.0.0.1:%s", envy.Get("PORT", "3000"))))
	return opts
}
