package servers

import (
	"context"
	"fmt"
	"net/http"
)

// Simple server
type Simple struct {
	*http.Server
}

// SetAddr sets the servers address, if it hasn't already been set
func (s *Simple) SetAddr(addr string) {
	if s.Server.Addr == "" {
		s.Server.Addr = addr
	}
}

// String returns a string representation of a Simple server
func (s *Simple) String() string {
	return fmt.Sprintf("simple server on %s", s.Server.Addr)
}

// Start the server
func (s *Simple) Start(c context.Context, h http.Handler) error {
	s.Handler = h
	return s.ListenAndServe()
}

// New Simple server
func New() *Simple {
	return &Simple{
		Server: &http.Server{},
	}
}
