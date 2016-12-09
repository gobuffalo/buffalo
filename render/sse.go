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

func (es *EventSource) Flush() {
	es.fl.Flush()
}

func (es *EventSource) CloseNotify() <-chan bool {
	return es.w.(http.CloseNotifier).CloseNotify()
}

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
