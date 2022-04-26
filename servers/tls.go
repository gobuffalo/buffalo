package servers

import (
	"context"
	"fmt"
	"net/http"
)

// TLS server
type TLS struct {
	*http.Server
	CertFile string
	KeyFile  string
}

// SetAddr sets the servers address, if it hasn't already been set
func (s *TLS) SetAddr(addr string) {
	if s.Server.Addr == "" {
		s.Server.Addr = addr
	}
}

// String returns a string representation of a Listener
func (s *TLS) String() string {
	return fmt.Sprintf("TLS server on %s", s.Server.Addr)
}

// Start the server
func (s *TLS) Start(c context.Context, h http.Handler) error {
	s.Handler = h
	return s.ListenAndServeTLS(s.CertFile, s.KeyFile)
}
