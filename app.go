package buffalo

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/gobuffalo/envy"
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

// Serve the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully. This will also start the Worker process, unless WorkerOff is enabled.
func (a *App) Serve() error {
	logrus.Infof("Starting application at %s", a.Options.Addr)
	server := http.Server{
		Handler: a,
	}
	ctx, cancel := sigtx.WithCancel(a.Context, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	go func() {
		// gracefully shut down the application when the context is cancelled
		<-ctx.Done()
		logrus.Info("Shutting down application")

		err := a.Stop(ctx.Err())
		if err != nil {
			logrus.Error(err)
		}

		if !a.WorkerOff {
			// stop the workers
			logrus.Info("Shutting down worker")
			err = a.Worker.Stop()
			if err != nil {
				logrus.Error(err)
			}
		}

		err = server.Shutdown(ctx)
		if err != nil {
			logrus.Error(err)
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

	var err error

	if strings.HasPrefix(a.Options.Addr, "unix:") {
		// Use an UNIX socket
		listener, err := net.Listen("unix", a.Options.Addr[5:])
		if err != nil {
			return a.Stop(err)
		}
		// start the web server
		err = server.Serve(listener)
	} else {
		// Use a TCP socket
		server.Addr = a.Options.Addr

		// start the web server
		err = server.ListenAndServe()
	}

	if err != nil {
		return a.Stop(err)
	}

	return nil
}

// Stop the application and attempt to gracefully shutdown
func (a *App) Stop(err error) error {
	a.cancel()
	if err != nil && errors.Cause(err) != context.Canceled {
		logrus.Error(err)
		return err
	}
	return nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws := &Response{
		ResponseWriter: w,
	}
	if a.MethodOverride != nil {
		a.MethodOverride(w, r)
	}
	if ok := a.processPreHandlers(ws, r); !ok {
		return
	}

	var h http.Handler
	h = a.router
	if a.Env == "development" {
		h = web.ErrorChecker(h)
	}
	h.ServeHTTP(ws, r)
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
	a.router.NotFoundHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		c := a.newContext(RouteInfo{}, res, req)
		err := errors.Errorf("path not found: %s %s", req.Method, req.URL.Path)
		a.ErrorHandlers.Get(404)(404, err, c)
	})

	if a.MethodOverride == nil {
		a.MethodOverride = MethodOverride
	}
	a.Use(a.PanicHandler)
	a.Use(RequestLogger)
	a.Use(sessionSaver)

	return a
}

func (a *App) processPreHandlers(res http.ResponseWriter, req *http.Request) bool {
	sh := func(h http.Handler) bool {
		h.ServeHTTP(res, req)
		if br, ok := res.(*Response); ok {
			if br.Status > 0 || br.Size > 0 {
				return false
			}
		}
		return true
	}

	for _, ph := range a.PreHandlers {
		if ok := sh(ph); !ok {
			return false
		}
	}

	last := http.Handler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {}))
	for _, ph := range a.PreWares {
		last = ph(last)
		if ok := sh(last); !ok {
			return false
		}
	}
	return true
}
