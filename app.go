package buffalo

import (
	"net/http"
	"sync"

	gcontext "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

// App is where it all happens! It holds on to options,
// the underlying router, the middleware, and more.
// Without an App you can't do much!
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

// New returns a new instance of App, without any frills
// or thrills. Most people will want to use Standard which
// adds some sane, and useful, defaults.
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

// Standard returns a new instace of App with sane defaults,
// some not so sane defaults, and a few bits and pieces to make
// your life that much easier. You'll want to use this almost
// all of the time to build your applications.
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
