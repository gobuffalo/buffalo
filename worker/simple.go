package worker

import (
	"time"
)

var _ Worker = Simple{}

type Simple struct{}

func (w Simple) Perform(job Job) error {
	if err := job.Valid(); err != nil {
		return err
	}
	go job.Handler(job.Args...)
	return nil
}

func (w Simple) PerformAt(job Job, t time.Time) error {
	return w.PerformIn(job, t.Sub(time.Now()))
}

func (w Simple) PerformIn(job Job, d time.Duration) error {
	if err := job.Valid(); err != nil {
		return err
	}
	time.AfterFunc(d, func() {
		w.Perform(job)
	})
	return nil
}
