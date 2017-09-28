package buffalo

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func voidHandler(c Context) error {
	return nil
}

func getDynamicPort(a *App) string {
	a.moot.Lock()
	defer a.moot.Unlock()
	_, port, _ := net.SplitHostPort(a.server.Addr)
	return port
}

func TestGracefulShutdown(t *testing.T) {
	requestFinished := false
	app := New(NewOptions())
	app.WorkerOff = true
	app.GET("/slow", func(c Context) error {
		time.Sleep(1 * time.Second)
		requestFinished = true
		return nil
	})

	go func() {
		time.Sleep(1500 * time.Millisecond)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	go func() {
		time.Sleep(750 * time.Millisecond)
		port := getDynamicPort(app)
		_, err := http.Get(fmt.Sprintf("http://localhost:%s/slow", port))
		if err != nil {
			t.Error("unexpected error requesting slow endpoint", err)
		}
	}()

	app.Start("127.0.0.1:0")

	require.True(t, requestFinished, "expected request to finish but was terminated early")
}

func TestGracefulShutdown_ForcedShutdown(t *testing.T) {
	requestFinished := false

	app := New(NewOptions())
	app.WorkerOff = true
	app.ShutDownTimeoutSeconds = 1
	app.GET("/slow", func(c Context) error {
		time.Sleep(4 * time.Second)
		requestFinished = true
		return nil
	})

	go func() {
		time.Sleep(1 * time.Second)
		app.Stop(nil)
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
		port := getDynamicPort(app)
		_, err := http.Get(fmt.Sprintf("http://localhost:%s/slow", port))
		if err != nil {
			t.Error("unexpected error requesting slow endpoint", err)
		}
	}()

	app.Start("127.0.0.1:0")

	require.False(t, requestFinished, "expected request to be terminated early")
}
