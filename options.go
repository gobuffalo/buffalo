package buffalo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/envy"
	"github.com/gorilla/sessions"
	"github.com/markbates/going/defaults"
)

// Options are used to configure and define how your application should run.
type Options struct {
	Name string
	// Env is the "environment" in which the App is running. Default is "development".
	Env string
	// LogLevel defaults to "debug".
	LogLevel string
	// Logger to be used with the application. A default one is provided.
	Logger Logger
	// MethodOverride allows for changing of the request method type. See the default
	// implementation at buffalo.MethodOverride
	MethodOverride http.HandlerFunc
	// SessionStore is the `github.com/gorilla/sessions` store used to back
	// the session. It defaults to use a cookie store and the ENV variable
	// `SESSION_SECRET`.
	SessionStore sessions.Store
	// SessionName is the name of the session cookie that is set. This defaults
	// to "_buffalo_session".
	SessionName string
	// Host that this application will be available at. Default is "http://127.0.0.1:[$PORT|3000]".
	Host string
	// Worker implements the Worker interface and can process tasks in the background.
	// Default is "github.com/gobuffalo/worker.Simple.
	Worker worker.Worker
	// WorkerOff tells App.Start() whether to start the Worker process or not. Default is "false".
	WorkerOff bool

	Context context.Context
	cancel  context.CancelFunc
	prefix  string
}

// NewOptions returns a new Options instance with sensible defaults
func NewOptions() Options {
	return optionsWithDefaults(Options{})
}

func optionsWithDefaults(opts Options) Options {
	opts.Env = defaults.String(opts.Env, envy.Get("GO_ENV", "development"))
	opts.LogLevel = defaults.String(opts.LogLevel, "debug")
	opts.Name = defaults.String(opts.Name, "/")

	if opts.Context == nil {
		opts.Context = context.Background()
	}
	opts.Context, opts.cancel = context.WithCancel(opts.Context)

	if opts.Logger == nil {
		opts.Logger = NewLogger(opts.LogLevel)
	}

	if opts.SessionStore == nil {
		secret := envy.Get("SESSION_SECRET", "")
		// In production a SESSION_SECRET must be set!
		if opts.Env == "production" && secret == "" {
			log.Println("WARNING! Unless you set SESSION_SECRET env variable, your session storage is not protected!")
		}
		opts.SessionStore = sessions.NewCookieStore([]byte(secret))
	}
	if opts.Worker == nil {
		opts.Worker = worker.NewSimpleWithContext(opts.Context)
	}
	opts.SessionName = defaults.String(opts.SessionName, "_buffalo_session")
	opts.Host = defaults.String(opts.Host, envy.Get("HOST", fmt.Sprintf("http://127.0.0.1:%s", envy.Get("PORT", "3000"))))
	return opts
}
