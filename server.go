package buffalo

import (
	"context"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/gobuffalo/buffalo/servers"
	"github.com/markbates/refresh/refresh/web"
	"github.com/markbates/sigtx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Serve the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully. This will also start the Worker process, unless WorkerOff is enabled.
func (a *App) Serve(srvs ...servers.Server) error {
	logrus.Infof("Starting application at %s", a.Options.Host)

	if len(srvs) == 0 {
		if strings.HasPrefix(a.Options.Addr, "unix:") {
			tcp, err := servers.UnixSocket(a.Options.Addr[5:])
			if err != nil {
				return errors.WithStack(err)
			}
			srvs = append(srvs, tcp)
		} else {
			srvs = append(srvs, servers.New())
		}
	}

	ctx, cancel := sigtx.WithCancel(a.Context, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	go func() {
		// gracefully shut down the application when the context is cancelled
		<-ctx.Done()
		logrus.Info("Shutting down application")

		if err := a.Stop(ctx.Err()); err != nil {
			logrus.Error(err)
		}

		if !a.WorkerOff {
			// stop the workers
			logrus.Info("Shutting down worker")
			if err := a.Worker.Stop(); err != nil {
				logrus.Error(err)
			}
		}

		for _, s := range srvs {
			if err := s.Shutdown(ctx); err != nil {
				logrus.Error(err)
			}
		}

	}()

	// if configured to do so, start the workers
	if !a.WorkerOff {
		go func() {
			if err := a.Worker.Start(ctx); err != nil {
				a.Stop(err)
			}
		}()
	}

	for _, s := range srvs {
		s.SetAddr(a.Addr)
		go func(s servers.Server) {
			if err := s.Start(ctx, a); err != nil {
				a.Stop(err)
			}
		}(s)
	}

	<-ctx.Done()
	return a.Context.Err()
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
