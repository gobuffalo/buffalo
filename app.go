package buffalo

import (
	"net/http"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"

	gcontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/markbates/going/defaults"
)

type App struct {
	Options
	router          *httprouter.Router
	moot            *sync.Mutex
	routes          routes
	root            *App
	middlewareStack middlewareStack
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer gcontext.Clear(r)
	ws := &buffaloResponse{
		ResponseWriter: w,
	}
	if a.MethodOverride != nil {
		a.MethodOverride(r)
	}
	a.router.ServeHTTP(ws, r)
}

func New(opts Options) *App {
	opts.Env = defaults.String(opts.Env, defaults.String(os.Getenv("BUFFALO_ENV"), defaults.String(os.Getenv("GO_ENV"), "development")))
	if opts.Logger == nil {
		l := log.New()
		l.Level, _ = log.ParseLevel(defaults.String(opts.LogLevel, "debug"))
		opts.Logger = l
	}
	if opts.MethodOverride == nil {
		opts.MethodOverride = MethodOverride
	}
	a := &App{
		Options:         opts,
		router:          httprouter.New(),
		moot:            &sync.Mutex{},
		routes:          routes{},
		middlewareStack: newMiddlewareStack(),
	}

	if a.NotFound == nil {
		a.NotFound = a.notFound()
	}
	a.router.NotFound = a.NotFound

	return a
}
