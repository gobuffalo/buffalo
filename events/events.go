package events

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
