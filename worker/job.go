package worker

import "errors"

// Job to be processed by a Worker
type Job struct {
	// Queue the job should be placed into
	Queue string
	// Args that will be passed to the Handler when run
	Args []interface{}
	// Handler that will be run by the worker
	Handler Handler
}

// Valid is this job a valid job?
func (j Job) Valid() error {
	if j.Handler == nil {
		return errors.New("must specify a Handler")
	}
	return nil
}
