package render

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type sse struct {
	Data interface{} `json:"data"`
	Type string      `json:"type"`
}

func (s *sse) String() string {
	b, _ := json.Marshal(s)
	return fmt.Sprintf("data: %s\n\n", string(b))
}

func (s *sse) Bytes() []byte {
	return []byte(s.String())
}

// EventSource is designed to work with JavaScript EventSource objects.
// see https://developer.mozilla.org/en-US/docs/Web/API/EventSource for
// more details
type EventSource struct {
	w  http.ResponseWriter
	fl http.Flusher
}

func (es *EventSource) Write(t string, d interface{}) error {
	s := &sse{Type: t, Data: d}
	_, err := es.w.Write(s.Bytes())
	if err != nil {
		return err
	}
	es.Flush()
	return nil
}

// Flush messages down the pipe. If messages aren't flushed they
// won't be sent.
func (es *EventSource) Flush() {
	es.fl.Flush()
}

// CloseNotify return true across the channel when the connection
// in the browser has been severed.
func (es *EventSource) CloseNotify() <-chan bool {
	return es.w.(http.CloseNotifier).CloseNotify()
}

// NewEventSource returns a new EventSource instance while ensuring
// that the http.ResponseWriter is able to handle EventSource messages.
// It also makes sure to set the proper response heads.
func NewEventSource(w http.ResponseWriter) (*EventSource, error) {
	es := &EventSource{w: w}
	var ok bool
	es.fl, ok = w.(http.Flusher)
	if !ok {
		return es, errors.New("Streaming is not supported!!")
	}

	es.w.Header().Set("Content-Type", "text/event-stream")
	es.w.Header().Set("Cache-Control", "no-cache")
	es.w.Header().Set("Connection", "keep-alive")
	es.w.Header().Set("Access-Control-Allow-Origin", "*")
	return es, nil
}
