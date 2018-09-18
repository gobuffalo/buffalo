package events

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/pkg/errors"
)

const (
	// AppStart is emitted when buffalo.App#Serve is called
	AppStart = "app:start"
	// AppStop is emitted when buffalo.App#Stop is called
	AppStop = "app:stop"
	// WorkerStart is emitted when buffalo.App#Serve is called and workers are started
	WorkerStart = "worker:start"
	// WorkerStop is emitted when buffalo.App#Stop is called and workers are stopped
	WorkerStop = "worker:stop"
	// RouteStarted is emitted when a requested route is being processed
	RouteStarted = "route:started"
	// RouteFinished is emitted when a requested route is completed
	RouteFinished = "route:finished"
	// ErrRoute is emitted when there is a problem handling processing a route
	ErrRoute = "err:route"
	// ErrGeneral is emitted for general errors
	ErrGeneral = "err:general"
	// ErrPanic is emitted when a panic is recovered
	ErrPanic = "err:panic"
	// ErrAppStart is emitted when an error occurs calling buffalo.App#Serve
	ErrAppStart = "err:app:start"
	// ErrAppStop is emitted when an error occurs calling buffalo.App#Stop
	ErrAppStop = "err:app:stop"
	// ErrWorkerStart is emitted when an error occurs when starting workers
	ErrWorkerStart = "err:worker:start"
	// ErrWorkerStop is emitted when an error occurs when stopping workers
	ErrWorkerStop = "err:worker:stop"
)

// Emit an event to all listeners
func Emit(e Event) error {
	return boss.Emit(e)
}

// Listen for events. Name is the name of the
// listener NOT the events you want to listen for
func Listen(name string, l Listener) {
	boss.Listen(name, l)
}

// StopListening removes the listener with the given name
func StopListening(name string) {
	boss.StopListening(name)
}

// SetManager allows you to replace the default
// event manager with a custom one
func SetManager(m Manager) {
	boss = m
}

// LoadPlugins will add listeners for any plugins that support "events"
func LoadPlugins() error {
	plugs, err := plugins.Available()
	if err != nil {
		return errors.WithStack(err)
	}
	for _, cmds := range plugs {
		for _, c := range cmds {
			if c.BuffaloCommand != "events" {
				continue
			}
			Listen(fmt.Sprintf("plugin-%s-%s", c.Binary, c.Name), func(e Event) {
				b, err := json.Marshal(e)
				if err != nil {
					fmt.Println("error trying to marshal event", e, err)
					return
				}
				cmd := exec.Command(c.Binary, c.UseCommand, string(b))
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				cmd.Stdin = os.Stdin
				if err := cmd.Run(); err != nil {
					fmt.Println("error trying to send event", strings.Join(cmd.Args, " "), err)
				}
			})
		}

	}
	return nil
}

func init() {
	LoadPlugins()
}
