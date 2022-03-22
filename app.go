package buffalo

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gobuffalo/envy"
	"github.com/gorilla/mux"
)

// App is where it all happens! It holds on to options,
// the underlying router, the middleware, and more.
// Without an App you can't do much!
type App struct {
	Options
	// Middleware, ErrorHandlers, router, and filepaths are moved to Home.
	Home
	moot   *sync.RWMutex
	routes RouteList
	// TODO: to be deprecated #road-to-v1
	root     *App
	children []*App

	// Routenamer for the app. This field provides the ability to override the
	// base route namer for something more specific to the app.
	RouteNamer RouteNamer
}

// Muxer returns the underlying mux router to allow
// for advance configurations
func (a *App) Muxer() *mux.Router {
	return a.router
}

// New returns a new instance of App and adds some sane, and useful, defaults.
func New(opts Options) *App {
	LoadPlugins()
	envy.Load()

	opts = optionsWithDefaults(opts)

	a := &App{
		Options: opts,
		Home: Home{
			name:   opts.Name,
			host:   opts.Host,
			prefix: opts.Prefix,
			ErrorHandlers: ErrorHandlers{
				http.StatusNotFound:            defaultErrorHandler,
				http.StatusInternalServerError: defaultErrorHandler,
			},
			router: mux.NewRouter(),
		},
		moot:     &sync.RWMutex{},
		routes:   RouteList{},
		children: []*App{},

		RouteNamer: baseRouteNamer{},
	}

	dem := a.defaultErrorMiddleware
	a.Middleware = newMiddlewareStack(dem)

	notFoundHandler := func(errorf string, code int) http.HandlerFunc {
		return func(res http.ResponseWriter, req *http.Request) {
			c := a.newContext(RouteInfo{}, res, req)
			err := fmt.Errorf(errorf, req.Method, req.URL.Path)
			_ = a.ErrorHandlers.Get(code)(code, err, c)
		}
	}

	a.router.NotFoundHandler = notFoundHandler("path not found: %s %s", http.StatusNotFound)
	a.router.MethodNotAllowedHandler = notFoundHandler("method not found: %s %s", http.StatusMethodNotAllowed)

	if a.MethodOverride == nil {
		a.MethodOverride = MethodOverride
	}
	a.Use(a.PanicHandler)
	a.Use(RequestLogger)

	return a
}
