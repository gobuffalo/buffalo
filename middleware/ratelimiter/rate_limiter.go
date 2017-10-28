package ratelimiter

import (
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
)

// Counter limits the amount of times an IP can run
type Counter interface {
	Increment(ip string) (count int, err error)
	Decrement(ip string) (count int, err error)
	Count(ip string) (count int, err error)
	Set(ip string, count int) error
}

// DefaultOptions is the default configuration for Rate Limiter middleware.
// By default, it fetches the IP based on IPHeaders
var DefaultOptions = Options{
	Operations: 100,
	Duration:   time.Second,
	GetIP: func(c buffalo.Context) (ip string, ok bool) {
		req := c.Request()
		for _, h := range []string{"X-Real-IP", "X-Client-IP", "X-Forwarded-For", "RemoteAddr"} {
			ip := req.Header.Get(h)
			if len(ip) > 0 {
				return ip, true
			}
		}
		return
	},
}

// Options define options for rate-limiter middleware
type Options struct {
	// Operations is how many operations per Duration an IP is allowed
	Operations int

	// Duration is how long to count Operations
	Duration time.Duration

	// GetIP is the func used to get the IP address from the request
	GetIP func(c buffalo.Context) (ip string, ok bool)
}

// Middleware is the http middleware that limits the interactions to the API
// based on the IP address using the provided Counter
func Middleware(r Counter, opts *Options) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			ip, ok := opts.GetIP(c)
			if !ok {
				c.Logger().Warn("could not detect IP")
				return next(c)
			}

			go func(ip string, opts *Options) {
				<-time.After(opts.Duration)
				r.Decrement(ip)
			}(ip, opts)

			v, err := r.Increment(ip)
			if err != nil {
				c.Logger().Errorf("could not increment counter for ip %q: %v", ip, err)
				return next(c)
			}

			if v <= opts.Operations {
				return next(c)
			}

			c.Logger().WithField("IP", ip).Warn("rate limited")
			c.Response().WriteHeader(http.StatusTooManyRequests)
			return nil
		}
	}
}
