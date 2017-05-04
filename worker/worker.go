package worker

import (
	"time"
)

// Handler function that will be run by the worker and given
// a slice of arguments
type Handler func(...interface{}) error

// Worker interface that needs to be implemented to be considered
// a "worker"
type Worker interface {
	// Perform a job as soon as possibly
	Perform(Job) error
	// PerformAt performs a job at a particular time
	PerformAt(Job, time.Time) error
	// PerformIn performs a job after waiting for a specified amount of time
	PerformIn(Job, time.Duration) error
}
