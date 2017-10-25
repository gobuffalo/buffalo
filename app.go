package buffalo

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

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
	children      []*App
	server        http.Server
	closed        bool
}

// Start the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully. This will also start the Worker process, unless WorkerOff is enabled.
func (a *App) Start(addr string) error {
	if !strings.Contains(addr, ":") {
		addr = fmt.Sprintf(":%s", addr)
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		a.Stop(err)
		return err
	}

	fmt.Printf("Starting application at %s:%d\n",
		listener.Addr().(*net.TCPAddr).IP,
		listener.Addr().(*net.TCPAddr).Port)

	a.moot.Lock()
	a.server = http.Server{
		Addr:    listener.Addr().String(),
		Handler: a,
	}
	a.moot.Unlock()

	errChan := make(chan error, 2)

	// if configured to do so, start the workers
	if !a.WorkerOff {
		go func() {
			err := a.Worker.Start(a.Context)
			if err != nil {
				errChan <- err
			}
		}()
	}

	go func() {
		errChan <- a.server.Serve(listener)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, os.Interrupt)

	for {
		select {
		case err := <-errChan:
			return a.Stop(err)
		case <-signalChan:
			return a.Stop(nil)
		}
	}
}

// Stop the application and attempt to gracefully shutdown
func (a *App) Stop(err error) error {
	a.moot.Lock()
	defer a.moot.Unlock()
	if !a.closed {
		fmt.Println("Shutting down application")
		a.closed = true
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(a.ShutDownTimeoutSeconds)*time.Second)
		defer cancel()
		a.server.Shutdown(shutdownCtx)

		if !a.WorkerOff {
			// stop the workers
			fmt.Println("Shutting down worker")
			err = a.Worker.Stop()
			if err != nil {
				fmt.Println(err)
			}
		}

		a.cancel()
		if err != nil && errors.Cause(err) != context.Canceled {
			fmt.Println(err)
			return err
		}
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

	if a.MethodOverride == nil {
		a.MethodOverride = MethodOverride
	}
	a.Use(a.PanicHandler)
	a.Use(RequestLogger)
	a.Use(sessionSaver)

	return a
}

// Automatic is deprecated, and will be removed in v0.10.0. Use buffalo.New instead.
func Automatic(opts Options) *App {
	warningMsg := "Automatic is deprecated, and will be removed in v0.10.0. Use buffalo.New instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}

	log.Println(warningMsg)
	return New(opts)
}

func (a *App) processPreHandlers(res http.ResponseWriter, req *http.Request) bool {
	sh := func(h http.Handler) bool {
		h.ServeHTTP(res, req)
		if br, ok := res.(*Response); ok {
			if (br.Status < 200 || br.Status > 299) && br.Status > 0 {
				return false
			}
			if br.Size > 0 {
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
