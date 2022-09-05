package worker

import (
	"context"
	"time"
)

// Handler function that will be run by the worker and given
// a slice of arguments
type Handler func(Args) error

// Worker interface that needs to be implemented to be considered
// a "worker"
type Worker interface {
	// Start the worker with the given context
	Start(context.Context) error
	// Stop the worker
	Stop() error
	// Perform a job as soon as possibly
	Perform(Job) error
	// PerformAt performs a job at a particular time
	PerformAt(Job, time.Time) error
	// PerformIn performs a job after waiting for a specified amount of time
	PerformIn(Job, time.Duration) error
	// Register a Handler
	Register(string, Handler) error
	// RegisterPeriodic performs a job periodically according to the provided cron spec
	RegisterPeriodic(cronSpec, jobName string, h Handler) error
}

/* TODO(sio4): #road-to-v1 - redefine Worker interface clearer
1. The Start() functions of current implementations including Simple,
   Gocraft Work Adapter do not block and immediately return the error.
   However, App.Serve() calls them within a go routine.
2. The Perform() family of functions can be called before the worker
   was started once the worker configured. Could be fine but there should
   be some guidiance for its usage.
3. The Perform() function could be interpreted as "Do it" by its name but
   their actual job is "Enqueue it" even though Simple worker has no clear
   boundary between them. It could make confusion.
*/
