package servers

import (
	"context"
	"net"
	"net/http"
)

// Server allows for custom server implementations
type Server interface {
	Addr() string
	Shutdown(context.Context) error
	Start(context.Context, http.Handler) error
	SetAddr(string)
}

// Wrap converts a standard *http.Server to a buffalo.Server
func Wrap(s *http.Server) Server {
	return &Simple{Server: s}
}

// WrapTLS Server converts a standard *http.Server to a buffalo.Server
// but makes sure it is run with TLS.
func WrapTLS(s *http.Server, certFile string, keyFile string) Server {
	return &TLS{
		Server:   s,
		CertFile: certFile,
		KeyFile:  keyFile,
	}
}

// WrapListener wraps an *http.Server and a net.Listener
func WrapListener(s *http.Server, l net.Listener) Server {
	return &Listener{
		Server:   s,
		Listener: l,
	}
}
