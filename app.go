package buffalo

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// App is where it all happens! It holds on to options,
// the underlying router, the middleware, and more.
// Without an App you can't do much!
type App struct {
	Options
	// Middleware returns the current MiddlewareStack for the App/Group.
	Middleware    *MiddlewareStack `json:"-"`
	ErrorHandlers ErrorHandlers    `json:"-"`
	router        *mux.Router
	moot          *sync.RWMutex
	routes        RouteList
	root          *App
	children      []*App
	filepaths     []string
}

// Muxer returns the underlying mux router to allow
// for advance configurations
func (a *App) Muxer() *mux.Router {
	return a.router
}

// New returns a new instance of App and adds some sane, and useful, defaults.
func New(opts Options) *App {
	loadEnv()

	(&opts).SensibleDefaults()

	a := &App{
		Options: opts,
		ErrorHandlers: ErrorHandlers{
			http.StatusNotFound:            defaultErrorHandler,
			http.StatusInternalServerError: defaultErrorHandler,
		},
		router:   mux.NewRouter(),
		moot:     &sync.RWMutex{},
		routes:   RouteList{},
		children: []*App{},
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
	a.Use(sessionSaver)

	return a
}

// Load .env files. Files will be loaded in the same order that are received.
// Redefined vars will override previously existing values.
// If no arg passed, it will try to load a .env file.
func loadEnv(files ...string) error {

	// If no files received, load the default one
	if len(files) == 0 {
		return godotenv.Overload()
	}

	// We received a list of files
	for _, file := range files {

		// Check if it exists or we can access
		if _, err := os.Stat(file); err != nil {
			// It does not exist or we can not access.
			// Return and stop loading
			return err
		}

		// It exists and we have permission. Load it
		if err := godotenv.Overload(file); err != nil {
			return err
		}

	}
	return nil
}
