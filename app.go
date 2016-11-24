package buffalo

import (
	"net/http"
	"sync"

	gcontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
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
	opts = optionsWithDefaults(opts)

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

func Standard() *App {
	a := New(NewOptions())
	if a.MethodOverride == nil {
		a.MethodOverride = MethodOverride
	}

	a.Use(RequestLogger)
	if a.NotFound == nil {
		a.NotFound = a.notFound()
	}
	a.router.NotFound = a.NotFound

	return a
}
