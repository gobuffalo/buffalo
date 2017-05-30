package buffalo

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

// Start the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully. This will also start the Worker process, unless WorkerOff is enabled.
func (a *App) Start(addr string) error {
	fmt.Printf("Starting application at %s\n", addr)
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: a,
	}

	go func() {
		<-a.Context.Done()
		fmt.Println("Shutting down application")
		a.cancel()
		err := server.Shutdown(a.Context)
		if err != nil {
			a.Logger.Error(errors.WithStack(err))
		}
		if !a.WorkerOff {
			err = a.Worker.Stop()
			if err != nil {
				a.Logger.Error(errors.WithStack(err))
			}
		}
	}()

	if !a.WorkerOff {
		go func() {
			err := a.Worker.Start(a.Context)
			if err != nil {
				a.Logger.Error(errors.WithStack(err))
				a.cancel()
			}
		}()
	}

	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		a.cancel()
	}()

	err := server.ListenAndServe()
	if err != nil {
		a.cancel()

		err = errors.WithStack(err)
		a.Logger.Error(err)
		return errors.WithStack(err)
	}
	return nil
}

// Stop the application and attempt to gracefully shutdown
func (a *App) Stop() error {
	a.cancel()
	return nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer gcontext.Clear(r)
	ws := &Response{
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
		Middleware: newMiddlewareStack(),
		ErrorHandlers: ErrorHandlers{
			404: defaultErrorHandler,
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
	a.Use(a.PanicHandler)
	a.Use(RequestLogger)

	return a
}
