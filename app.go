package buffalo

import (
	"net/http"
	"sync"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/markbates/refresh/refresh/web"
	"github.com/pkg/errors"
)

// App is where it all happens! It holds on to options,
// the underlying router, the middleware, and more.
// Without an App you can't do much!
type App struct {
	Options
	// Middleware returns the current MiddlewareStack for the App/Group.
	Middleware    *MiddlewareStack
	ErrorHandlers ErrorHandlers
	router        *mux.Router
	moot          *sync.Mutex
	routes        RouteList
	root          *App
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer gcontext.Clear(r)
	ws := &buffaloResponse{
		ResponseWriter: w,
	}
	if a.MethodOverride != nil {
		a.MethodOverride(w, r)
	}
	var h http.Handler
	h = a.router
	a.addDefaultRouteNames()

	if a.Env == "development" {
		h = web.ErrorChecker(h)
	}

	h.ServeHTTP(ws, r)
}

// New returns a new instance of App, without any frills
// or thrills. Most people will want to use Automatic which
// adds some sane, and useful, defaults.
func New(opts Options) *App {
	opts = optionsWithDefaults(opts)

	a := &App{
		Options:    opts,
		Middleware: newMiddlewareStack(),
		ErrorHandlers: ErrorHandlers{
			404: NotFoundHandler,
			500: defaultErrorHandler,
		},
		router: mux.NewRouter(),
		moot:   &sync.Mutex{},
		routes: RouteList{},
	}
	a.router.NotFoundHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		c := a.newContext(RouteInfo{}, res, req)
		err := errors.Errorf("path not found: %s", req.URL.Path)
		a.ErrorHandlers.Get(404)(404, err, c)
	})

	return a
}

// Automatic returns a new instance of App with sane defaults,
// some not so sane defaults, and a few bits and pieces to make
// your life that much easier. You'll want to use this almost
// all of the time to build your applications.
//
// https://www.youtube.com/watch?v=BKbOplYmjZM
func Automatic(opts Options) *App {
	opts = optionsWithDefaults(opts)

	a := New(opts)

	if a.MethodOverride == nil {
		a.MethodOverride = MethodOverride
	}

	a.Use(RequestLogger)

	return a
}

func (a *App) addDefaultRouteNames() {
	for _, route := range a.Routes() {
		if route.Name() == "" {
			route.RouteName(buildRouteName(route.Path))
		}
	}
}
