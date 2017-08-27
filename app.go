package buffalo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"

	gcontext "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/markbates/refresh/refresh/web"
	"github.com/markbates/sigtx"
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

// Start the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully. This will also start the Worker process, unless WorkerOff is enabled.
func (a *App) Start(addr string) error {
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf(":%s", addr)
	}
	fmt.Printf("Starting application at %s\n", addr)
	server := http.Server{
		Addr:    addr,
		Handler: a,
	}
	ctx, cancel := sigtx.WithCancel(a.Context, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	go func() {
		// gracefully shut down the application when the context is cancelled
		<-ctx.Done()
		fmt.Println("Shutting down application")

		err := a.Stop(ctx.Err())
		if err != nil {
			fmt.Println(err)
		}

		if !a.WorkerOff {
			// stop the workers
			fmt.Println("Shutting down worker")
			err = a.Worker.Stop()
			if err != nil {
				fmt.Println(err)
			}
		}

		err = server.Shutdown(ctx)
		if err != nil {
			fmt.Println(err)
		}

	}()

	// if configured to do so, start the workers
	if !a.WorkerOff {
		go func() {
			err := a.Worker.Start(ctx)
			if err != nil {
				a.Stop(err)
			}
		}()
	}

	// start the web server
	err := server.ListenAndServe()
	if err != nil {
		return a.Stop(err)
	}
	return nil
}

// Stop the application and attempt to gracefully shutdown
func (a *App) Stop(err error) error {
	a.cancel()
	if err != nil && errors.Cause(err) != context.Canceled {
		fmt.Println(err)
		return err
	}
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
		router:   mux.NewRouter().StrictSlash(true),
		moot:     &sync.Mutex{},
		routes:   RouteList{},
		children: []*App{},
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
