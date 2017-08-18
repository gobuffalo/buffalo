// +build appengine

package buffalo

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Start the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully. This will also start the Worker process, unless WorkerOff is enabled.
func (a *App) Start(addr string) error {
	http.Handle("/", a)

	go func() {
		// gracefully shut down the application when the context is cancelled
		<-a.Context.Done()
		fmt.Println("Shutting down application")
		if !a.WorkerOff {
			// stop the workers
			err := a.Worker.Stop()
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

	return nil
}
