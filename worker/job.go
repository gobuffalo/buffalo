package worker

import "errors"

type Job struct {
	Queue   string
	Args    []interface{}
	Handler Handler
}

func (j Job) Valid() error {
	if j.Handler == nil {
		return errors.New("must specify a Handler")
	}
	return nil
}
