package ratelimiter

import (
	"net/http"
	"testing"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/assert"
)

func memoryRateLimiterTestApp() (cnt Counter, app *buffalo.App) {
	app = buffalo.New(buffalo.Options{})
	cnt = NewMemoryCounter()
	opts := DefaultOptions
	opts.Duration = time.Millisecond * 10
	app.Use(Middleware(cnt, &opts))
	app.ANY("/", func(c buffalo.Context) error {
		c.Response().WriteHeader(http.StatusOK)
		return nil
	})
	return
}

func TestRateLimiterShouldNotLimitOnFirstRun(t *testing.T) {
	cnt, app := memoryRateLimiterTestApp()
	ip := "abcd"
	cnt.Set(ip, 1)

	w := willie.New(app)
	req := w.HTML("/")
	req.Headers["x-real-ip"] = ip
	resp := req.Get()

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestRateLimiterShouldLimitWhenCounterEqualToLimit(t *testing.T) {
	cnt, app := memoryRateLimiterTestApp()

	ip := "abcd"
	cnt.Set(ip, DefaultOptions.Operations)

	w := willie.New(app)
	req := w.HTML("/")
	req.Headers["x-real-ip"] = ip
	resp := req.Get()

	assert.Equal(t, http.StatusTooManyRequests, resp.Code)
}

func TestRateLimiterShouldLimitWhenCounterOveLimit(t *testing.T) {
	cnt, app := memoryRateLimiterTestApp()

	ip := "abcd"
	cnt.Set(ip, DefaultOptions.Operations+1)

	w := willie.New(app)
	req := w.HTML("/")
	req.Headers["x-real-ip"] = ip
	resp := req.Get()

	assert.Equal(t, http.StatusTooManyRequests, resp.Code)
}
func TestRateLimiterShouldNotLimitWhenUnknownIP(t *testing.T) {
	_, app := memoryRateLimiterTestApp()

	w := willie.New(app)
	req := w.HTML("/")
	resp := req.Get()
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestRateLimiterShouldIncrementCounter(t *testing.T) {
	cnt, app := memoryRateLimiterTestApp()
	ip := "abcd"
	cnt.Set(ip, 1)

	w := willie.New(app)
	req := w.HTML("/")
	req.Headers["x-real-ip"] = ip
	req.Get()

	n, _ := cnt.Count(ip)
	assert.Equal(t, 2, n)
}

func TestRateLimiterShouldDecrementCounter(t *testing.T) {
	cnt, app := memoryRateLimiterTestApp()
	ip := "abcd"
	cnt.Set(ip, 1)

	w := willie.New(app)
	req := w.HTML("/")
	req.Headers["x-real-ip"] = ip
	req.Get()
	time.Sleep(DefaultOptions.Duration)

	n, _ := cnt.Count(ip)
	assert.Equal(t, 1, n)
}
