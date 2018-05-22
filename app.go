package buffalo

import (
	"net/http"
	"sync"

	"github.com/gobuffalo/envy"
	"github.com/gorilla/mux"
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
	children      []*App
}

// New returns a new instance of App and adds some sane, and useful, defaults.
func New(opts Options) *App {
	envy.Load()
	opts = optionsWithDefaults(opts)

	a := &App{
		Options:    opts,
		Middleware: newMiddlewareStack(),
		ErrorHandlers: ErrorHandlers{
			404: defaultErrorHandler,
			500: defaultErrorHandler,
		},
		router:   mux.NewRouter().StrictSlash(!opts.LooseSlash),
		moot:     &sync.Mutex{},
		routes:   RouteList{},
		children: []*App{},
	}

	notFoundHandler := func(errorf string, code int) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			c := a.newContext(RouteInfo{}, res, req)
			err := errors.Errorf(errorf, req.Method, req.URL.Path)
			a.ErrorHandlers.Get(code)(code, err, c)
		}
	}

	a.router.NotFoundHandler = notFoundHandler("path not found: %s %s", 404)
	a.router.MethodNotAllowedHandler = notFoundHandler("method not found: %s %s", 405)

	if a.MethodOverride == nil {
		a.MethodOverride = MethodOverride
	}
	a.Use(a.PanicHandler)
	a.Use(RequestLogger)
	a.Use(sessionSaver)

	return a
}
