package servers

import (
	"context"
	"net"
	"net/http"
)

// Listener server for using a pre-defined net.Listener
type Listener struct {
	*http.Server
	Listener net.Listener
}

// SetAddr sets the servers address, if it hasn't already been set
func (s *Listener) SetAddr(addr string) {
	if s.Server.Addr == "" {
		s.Server.Addr = addr
	}
}

// Start the server
func (s *Listener) Start(c context.Context, h http.Handler) error {
	s.Handler = h
	return s.Serve(s.Listener)
}

// UnixSocket returns a new Listener on that address
func UnixSocket(addr string) (*Listener, error) {
	listener, err := net.Listen("unix", addr)
	if err != nil {
		return nil, err
	}
	return &Listener{
		Server:   &http.Server{},
		Listener: listener,
	}, nil
}
