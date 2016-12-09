package buffalo

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/markbates/refresh/refresh/web"
)

// App is where it all happens! It holds on to options,
// the underlying router, the middleware, and more.
// Without an App you can't do much!
type App struct {
	Options
	// Middleware returns the current MiddlewareStack for the App/Group.
	Middleware *MiddlewareStack
	router     *mux.Router
	moot       *sync.Mutex
	routes     RouteList
	root       *App
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
		router:     mux.NewRouter(),
		moot:       &sync.Mutex{},
		routes:     RouteList{},
		Middleware: newMiddlewareStack(),
	}
	if a.Logger == nil {
		a.Logger = NewLogger(opts.LogLevel)
	}
	if a.NotFound == nil {
		a.NotFound = a.notFound()
	}
	a.router.NotFoundHandler = a.NotFound

	return a
}

// Automatic returns a new instace of App with sane defaults,
// some not so sane defaults, and a few bits and pieces to make
// your life that much easier. You'll want to use this almost
// all of the time to build your applications.
//
// https://www.youtube.com/watch?v=BKbOplYmjZM
func Automatic(opts Options) *App {
	opts = optionsWithDefaults(opts)
	if opts.MethodOverride == nil {
		opts.MethodOverride = MethodOverrideFunc
	}
	if opts.Logger == nil {
		lvl, _ := logrus.ParseLevel(opts.LogLevel)

		hl := logrus.New()
		hl.Level = lvl
		hl.Formatter = &logrus.TextFormatter{}
		// hl.Out = os.Stdout

		err := os.MkdirAll(opts.LogDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(filepath.Join(opts.LogDir, opts.Env+".log"))
		if err != nil {
			log.Fatal(err)
		}
		fl := logrus.New()
		fl.Level = lvl
		fl.Formatter = &logrus.JSONFormatter{}
		fl.Out = f

		ml := &multiLogger{Loggers: []logrus.FieldLogger{hl, fl}}
		opts.Logger = ml
	}

	a := New(opts)

	if a.MethodOverride == nil {
		a.MethodOverride = MethodOverride
	}

	a.Use(RequestLogger)
	if a.NotFound == nil {
		a.NotFound = a.notFound()
	}
	a.router.NotFoundHandler = a.NotFound

	return a
}
