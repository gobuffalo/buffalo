package worker

import "encoding/json"

// Args are the arguments passed into a job
type Args map[string]interface{}

func (a Args) String() string {
	b, _ := json.Marshal(a)
	return string(b)
}

// Job to be processed by a Worker
type Job struct {
	// Queue the job should be placed into
	Queue string
	// Args that will be passed to the Handler when run
	Args Args
	// Handler that will be run by the worker
	Handler string
}

func (j Job) String() string {
	b, _ := json.Marshal(j)
	return string(b)
}
