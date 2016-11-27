package buffalo

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Sirupsen/logrus"
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
// or thrills. Most people will want to use Automatic which
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
	if a.Logger == nil {
		l := logrus.New()
		l.Level, _ = logrus.ParseLevel(opts.LogLevel)
		ml := &MultiLogger{Loggers: []logrus.FieldLogger{l}}
		a.Logger = ml
	}
	if a.NotFound == nil {
		a.NotFound = a.notFound()
	}
	a.router.NotFound = a.NotFound

	return a
}

// Automatic returns a new instace of App with sane defaults,
// some not so sane defaults, and a few bits and pieces to make
// your life that much easier. You'll want to use this almost
// all of the time to build your applications.
func Automatic(opts Options) *App {
	opts = optionsWithDefaults(opts)
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

		ml := &MultiLogger{Loggers: []logrus.FieldLogger{hl, fl}}
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
	a.router.NotFound = a.NotFound

	return a
}
