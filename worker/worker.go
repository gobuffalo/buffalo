package worker

import (
	"time"
)

type Handler func(...interface{}) error

type Worker interface {
	Perform(Job) error
	PerformAt(Job, time.Time) error
	PerformIn(Job, time.Duration) error
}
