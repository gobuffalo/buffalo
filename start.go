// +build !appengine

package buffalo

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
)

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
		// gracefully shut down the application when the context is cancelled
		<-a.Context.Done()
		fmt.Println("Shutting down application")
		err := server.Shutdown(a.Context)
		if err != nil {
			a.Logger.Error(errors.WithStack(err))
		}
		if !a.WorkerOff {
			// stop the workers
			err = a.Worker.Stop()
			if err != nil {
				a.Logger.Error(errors.WithStack(err))
			}
		}
	}()

	// if configured to do so, start the workers
	if !a.WorkerOff {
		go func() {
			err := a.Worker.Start(a.Context)
			if err != nil {
				a.Stop(errors.WithStack(err))
			}
		}()
	}

	// listen for system signals, like CTRL-C
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		<-signalChan
		a.Stop(nil)
	}()

	// start the web server
	err := server.ListenAndServe()
	if err != nil {
		return a.Stop(errors.WithStack(err))
	}
	return nil
}
