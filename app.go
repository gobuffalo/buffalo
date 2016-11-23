package buffalo

import (
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	humanize "github.com/flosch/go-humanize"
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
	now := time.Now()
	l := a.Logger.WithFields(log.Fields{
		"method": r.Method,
		"path":   r.URL,
	})

	ws := &buffaloResponse{
		ResponseWriter: w,
		logger:         l,
	}
	defer func() {
		l = ws.logger.WithFields(log.Fields{
			"duration": time.Now().Sub(now),
			"size":     humanize.Bytes(uint64(ws.size)),
			"status":   ws.status,
		})
		l.Info()
	}()
	a.router.ServeHTTP(ws, r)
}

func New(opts Options) *App {
	opts.Env = defaults.String(opts.Env, defaults.String(os.Getenv("BUFFALO_ENV"), defaults.String(os.Getenv("GO_ENV"), "development")))
	if opts.Logger == nil {
		l := log.New()
		l.Level, _ = log.ParseLevel(defaults.String(opts.LogLevel, "debug"))
		opts.Logger = l
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
