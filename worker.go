package buffalo

import "time"

type WorkerFunc func() error

type Worker interface {
	Perform(WorkerFunc) error
	PerformAt(WorkerFunc, time.Time) error
	PerformIn(WorkerFunc, time.Duration) error
}

var _ Worker = defaultWorker{}

type defaultWorker struct {
}

func (w defaultWorker) Perform(wf WorkerFunc) error {
	go wf()
	return nil
}

func (w defaultWorker) PerformAt(wf WorkerFunc, t time.Time) error {
	return w.PerformIn(wf, t.Sub(time.Now()))
}

func (w defaultWorker) PerformIn(wf WorkerFunc, d time.Duration) error {
	time.AfterFunc(d, func() {
		wf()
	})
	return nil
}
