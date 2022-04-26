package buffalo

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gobuffalo/buffalo/servers"
	"github.com/gobuffalo/events"
	"github.com/markbates/refresh/refresh/web"
)

// Serve the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully. This will also start the Worker process, unless WorkerOff is enabled.
func (a *App) Serve(srvs ...servers.Server) error {
	var wg sync.WaitGroup

	a.Logger.Debug("starting application")

	payload := events.Payload{
		"app": a,
	}
	if err := events.EmitPayload(EvtAppStart, payload); err != nil {
		// just to make sure if events work properly?
		a.Logger.Error("unable to emit event. something went wrong internally")
		return err
	}

	if len(srvs) == 0 {
		if strings.HasPrefix(a.Options.Addr, "unix:") {
			tcp, err := servers.UnixSocket(a.Options.Addr[5:])
			if err != nil {
				return err
			}
			srvs = append(srvs, tcp)
		} else {
			srvs = append(srvs, servers.New())
		}
	}

	ctx, cancel := signal.NotifyContext(a.Context, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	wg.Add(1)
	go func() {
		// gracefully shut down the application when the context is cancelled
		defer wg.Done()
		// channel waiter should not be called any other place
		<-ctx.Done()

		a.Logger.Info("shutting down application")

		// shutting down listeners first, to make sure no more new request
		a.Logger.Info("shutting down servers")
		for _, s := range srvs {
			timeout := time.Duration(a.Options.TimeoutSecondShutdown) * time.Second
			ctx, cfn := context.WithTimeout(context.Background(), timeout)
			defer cfn()
			events.EmitPayload(EvtServerStop, payload)
			if err := s.Shutdown(ctx); err != nil {
				events.EmitError(EvtServerStopErr, err, payload)
				a.Logger.Error("shutting down server: ", err)
			}
			cfn()
		}

		if !a.WorkerOff {
			a.Logger.Info("shutting down worker")
			events.EmitPayload(EvtWorkerStop, payload)
			if err := a.Worker.Stop(); err != nil {
				events.EmitError(EvtWorkerStopErr, err, payload)
				a.Logger.Error("error while shutting down worker: ", err)
			}
		}
	}()

	// if configured to do so, start the workers
	if !a.WorkerOff {
		wg.Add(1)
		go func() {
			defer wg.Done()
			events.EmitPayload(EvtWorkerStart, payload)
			if err := a.Worker.Start(ctx); err != nil {
				events.EmitError(EvtWorkerStartErr, err, payload)
				a.Stop(err)
			}
		}()
	}

	for _, s := range srvs {
		s.SetAddr(a.Addr)
		a.Logger.Infof("starting %s", s)
		wg.Add(1)
		go func(s servers.Server) {
			defer wg.Done()
			events.EmitPayload(EvtServerStart, payload)
			// s.Start always returns non-nil error
			a.Stop(s.Start(ctx, a))
		}(s)
	}

	wg.Wait()
	a.Logger.Info("shutdown completed")

	err := ctx.Err()
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

// Stop the application and attempt to gracefully shutdown
func (a *App) Stop(err error) error {
	events.EmitError(EvtAppStop, err, events.Payload{"app": a})

	ce := a.Context.Err()
	if ce != nil {
		a.Logger.Warn("application context has already been canceled: ", ce)
		return errors.New("application has already been canceled")
	}

	a.Logger.Warn("stopping application: ", err)
	a.cancel()
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

	r.URL.Path = a.normalizePath(r.URL.Path)

	var h http.Handler = a.router
	if a.Env == "development" {
		h = web.ErrorChecker(h)
	}
	h.ServeHTTP(ws, r)
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

func (a *App) normalizePath(path string) string {
	if strings.HasSuffix(path, "/") {
		return path
	}
	for _, p := range a.filepaths {
		if p == "/" {
			continue
		}
		if strings.HasPrefix(path, p) {
			return path
		}
	}
	return path + "/"
}
