package buffalo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gobuffalo/buffalo/internal/consts"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/logger"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/logging"
	"github.com/gorilla/sessions"
)

type Env string

func (e Env) String() string {
	if len(e) == 0 {
		return consts.Development
	}
	return string(e)
}

func (e Env) Development() bool {
	if len(e) == 0 {
		return true
	}
	return string(e) == consts.Development
}

func (e Env) Test() bool {
	return string(e) == consts.Test
}

func (e Env) Production() bool {
	return string(e) == consts.Production
}

func (e Env) NotProd() bool {
	return string(e) != consts.Production
}

// Options are used to configure and define how your application should run.
type Options struct {
	Name string `json:"name"`
	// Addr is the bind address provided to http.Server. Default is "127.0.0.1:3000"
	// Can be set using ENV vars "ADDR" and "PORT".
	Addr string `json:"addr"`
	// Host that this application will be available at. Default is "http://127.0.0.1:[$PORT|3000]".
	Host string `json:"host"`

	// Port that this application will be available at. Defaults is "3000. Can be set using ENV var "PORT"
	Port string `json:"port"`

	// Env is the "environment" in which the App is running. Default is "development".
	Env string `json:"env"`

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

// NewOptions returns a new Options instance with sensible defaults
func NewOptions() Options {
	var opts Options
	(&opts).SensibleDefaults()
	return opts
}

func (opts *Options) defEnv() error {
	if len(opts.Env) == 0 {
		opts.Env = os.Getenv(consts.GO_ENV)
	}

	if len(opts.Env) == 0 {
		opts.Env = consts.Development
	}
	return nil
}

func (opts *Options) defPort() error {
	if len(opts.Port) > 0 {
		return nil
	}

	port := os.Getenv(consts.PORT)
	if len(port) == 0 {
		port = consts.Def_Port
	}
	opts.Port = port
	return nil
}

func (opts *Options) defAddr() error {
	addr := consts.Def_Addr
	if Env(opts.Env).Development() {
		addr = consts.Def_AddrDev
	}
	envAddr := os.Getenv(consts.ADDR)
	if len(envAddr) == 0 {
		envAddr = addr
	}

	const unix = "unix:"
	if strings.HasPrefix(envAddr, unix) {
		// UNIX domain socket doesn't have a port
		opts.Addr = envAddr
	}
	if len(opts.Addr) > 0 {
		return nil
	}
	if err := opts.defPort(); err != nil {
		return err
	}
	opts.Addr = fmt.Sprintf("%s:%s", envAddr, opts.Port)
	return nil
}

func (opts *Options) defLogger() error {
	if opts.Logger != nil {
		return nil
	}
	var err error
	lvl := os.Getenv(consts.LOG_LEVEL)
	if len(lvl) > 0 {
		opts.LogLvl, err = logger.ParseLevel(lvl)
		if err != nil {
			opts.LogLvl = logger.DebugLevel
		}
	}
	if opts.LogLvl == 0 {
		opts.LogLvl = logger.DebugLevel
	}
	opts.Logger = logger.New(opts.LogLvl)

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
	return nil
}

func (opts *Options) defSession() error {
	if len(opts.SessionName) == 0 {
		opts.SessionName = consts.Def_SessionName
	}

	if opts.SessionStore != nil {
		return nil
	}
	secret := os.Getenv(consts.SESSION_SECRET)

	const bufsec = "buffalo-secret"
	env := Env(opts.Env)
	if len(secret) == 0 && (env.Development() || env.Test()) {
		secret = bufsec
	}

	// In production a SESSION_SECRET must be set!
	if len(secret) == 0 {
		opts.Logger.Warn("Unless you set SESSION_SECRET env variable, your session storage is not protected!")
	}

	cookieStore := sessions.NewCookieStore([]byte(secret))

	//Cookie secure attributes, see: https://www.owasp.org/index.php/Testing_for_cookies_attributes_(OTG-SESS-002)
	cookieStore.Options.HttpOnly = true
	if opts.Env == "production" {
		cookieStore.Options.Secure = true
	}

	opts.SessionStore = cookieStore
	return nil
}

func (opts *Options) defWorker() error {
	if opts.Worker != nil {
		return nil
	}
	w := worker.NewSimpleWithContext(opts.Context)
	w.Logger = opts.Logger
	opts.Worker = w
	return nil
}

func (opts *Options) defHost() error {
	if len(opts.Host) > 0 {
		return nil
	}

	host := os.Getenv(consts.HOST)
	if len(host) > 0 {
		opts.Host = host
		return nil
	}

	if err := opts.defPort(); err != nil {
		return err
	}

	opts.Host = fmt.Sprintf("http://%s:%s", opts.Addr, opts.Port)

	return nil
}

// SensibleDefaults will set any unset values to sensible defaults values.
func (opts *Options) SensibleDefaults() error {
	if err := opts.defEnv(); err != nil {
		return err
	}
	if err := opts.defAddr(); err != nil {
		return err
	}

	if len(opts.Name) == 0 {
		opts.Name = consts.Def_Root
	}

	if opts.Context == nil {
		opts.Context = context.Background()
	}

	opts.Context, opts.cancel = context.WithCancel(opts.Context)

	if err := opts.defLogger(); err != nil {
		return err
	}

	if err := opts.defSession(); err != nil {
		return err
	}

	if err := opts.defWorker(); err != nil {
		return err
	}

	if len(opts.Host) == 0 {
	}
	return nil
}
