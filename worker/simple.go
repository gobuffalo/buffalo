package worker

import (
	"time"
)

var _ Worker = Simple{}

// Simple is a basic implementation of the Worker interface
// that is backed using just the standard library and goroutines.
type Simple struct{}

// Perform a job as soon as possibly using a goroutine.
func (w Simple) Perform(job Job) error {
	if err := job.Valid(); err != nil {
		return err
	}
	go job.Handler(job.Args...)
	return nil
}

// PerformAt performs a job at a particular time using a goroutine.
func (w Simple) PerformAt(job Job, t time.Time) error {
	return w.PerformIn(job, t.Sub(time.Now()))
}

// PerformIn performs a job after waiting for a specified amount
// using a goroutine.
func (w Simple) PerformIn(job Job, d time.Duration) error {
	if err := job.Valid(); err != nil {
		return err
	}
	time.AfterFunc(d, func() {
		w.Perform(job)
	})
	return nil
}
