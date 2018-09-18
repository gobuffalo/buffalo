package events

import "errors"

// Event represents different events
// in the lifecycle of a Buffalo app
type Event struct {
	// Kind is the "type" of event "app:start"
	Kind string
	// Message is optional
	Message string
	// Payload is optional
	Payload interface{}
	// Error is optional
	Error error
}

// Validate that an event is ready to be emitted
func (e Event) Validate() error {
	if len(e.Kind) == 0 {
		return errors.New("kind can not be blank")
	}
	return nil
}
