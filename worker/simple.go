package worker

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var _ Worker = &simple{}

// NewSimple creates a basic implementation of the Worker interface
// that is backed using just the standard library and goroutines.
func NewSimple() Worker {
	return NewSimpleWithContext(context.Background())
}

// NewSimpleWithContext creates a basic implementation of the Worker interface
// that is backed using just the standard library and goroutines.
func NewSimpleWithContext(ctx context.Context) Worker {
	ctx, cancel := context.WithCancel(ctx)
	return &simple{
		ctx:      ctx,
		cancel:   cancel,
		handlers: map[string]Handler{},
		moot:     &sync.Mutex{},
	}
}

// simple is a basic implementation of the Worker interface
// that is backed using just the standard library and goroutines.
type simple struct {
	ctx      context.Context
	cancel   context.CancelFunc
	handlers map[string]Handler
	moot     *sync.Mutex
}

func (w *simple) Register(name string, h Handler) error {
	w.moot.Lock()
	defer w.moot.Unlock()
	if _, ok := w.handlers[name]; ok {
		return errors.Errorf("handler already mapped for name %s", name)
	}
	w.handlers[name] = h
	return nil
}

func (w *simple) Start(ctx context.Context) error {
	w.ctx, w.cancel = context.WithCancel(ctx)
	return nil
}

func (w simple) Stop() error {
	w.cancel()
	return nil
}

// Perform a job as soon as possibly using a goroutine.
func (w simple) Perform(job Job) error {
	w.moot.Lock()
	defer w.moot.Unlock()
	if h, ok := w.handlers[job.Handler]; ok {
		go h(job.Args)
		return nil
	}
	return errors.Errorf("no handler mapped for name %s", job.Handler)
}

// PerformAt performs a job at a particular time using a goroutine.
func (w simple) PerformAt(job Job, t time.Time) error {
	return w.PerformIn(job, t.Sub(time.Now()))
}

// PerformIn performs a job after waiting for a specified amount
// using a goroutine.
func (w simple) PerformIn(job Job, d time.Duration) error {
	go func() {
		select {
		case <-time.After(d):
			w.Perform(job)
		case <-w.ctx.Done():
			w.cancel()
		}
	}()
	return nil
}
