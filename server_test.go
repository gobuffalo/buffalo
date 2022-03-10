package buffalo

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/stretchr/testify/require"
)

// All tests in this file requires certain amount of waiting and they are
// timing sensitive. Adjust this timing values if they are failing due to
// timing issue.
const (
	WAIT_START   = 2
	WAIT_RUN     = 2
	CONSUMER_RUN = 8
)

// startApp starts given buffalo app and check its exit status.
// The go routine emulates a buffalo app process.
func startApp(app *App, wg *sync.WaitGroup, r *require.Assertions) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := app.Serve()
		r.NoError(err)
	}()
	// wait until the server started.
	// could be improved with connection test but that's too much...
	time.Sleep(WAIT_START * time.Second)
}

func Test_Server_Simple(t *testing.T) {
	// This testcase explains the minimum/basic workflow of buffalo app.
	// Setup and execute the app, wait until startup, then stop it.
	// The other testcases use this structure with additional actions.
	r := require.New(t)
	var wg sync.WaitGroup

	// Setup a new buffalo.App to be used as a testing buffalo app.
	app := New(Options{})

	startApp(app, &wg, r) // starts buffalo app routine.

	app.cancel()
	wg.Wait()
}

var handlerDone = false

// timeConsumer consumes about 10 minutes for processing its request
func timeConsumer(c Context) error {
	for i := 0; i < CONSUMER_RUN; i++ {
		fmt.Println("#")
		time.Sleep(1 * time.Second)
	}
	handlerDone = true
	return c.Render(http.StatusOK, render.String("Hey!"))
}

func Test_Server_GracefulShutdownOngoingRequest(t *testing.T) {
	// This test case explain the minimum/basic workflow of buffalo app.
	r := require.New(t)
	var wg sync.WaitGroup

	// Setup a new buffalo.App with a simple time consuming handler.
	app := New(Options{})
	app.GET("/", timeConsumer)

	startApp(app, &wg, r) // starts buffalo app routine.

	firstQuery := false
	secondQuery := false
	// This routine is the 1st client that GETs before Stop it
	// The result should be successful even though the server shutting down.
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := http.Get("http://127.0.0.1:3000")
		r.NoError(err)
		defer resp.Body.Close()
		r.Equal(http.StatusOK, resp.StatusCode)
		fmt.Println("the first query should be OK:", resp.Status)
		firstQuery = true
	}()
	// make sure the request sent
	time.Sleep(WAIT_RUN * time.Second)

	app.cancel()
	time.Sleep(1 * time.Second) // make sure the server started shutdown.

	// This routine is the 2nd client that GETs after Stop it
	// The result should be connection refused even though app is still on.
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := http.Get("http://127.0.0.1:3000")
		r.Contains(err.Error(), "refused")
		fmt.Println("the second query should be refused:", err)
		secondQuery = true
	}()

	wg.Wait()
	r.Equal(true, handlerDone)
	r.Equal(true, firstQuery)
	r.Equal(true, secondQuery)
}

var timerDone = false

func timerWorker(args worker.Args) error {
	for i := 0; i < CONSUMER_RUN; i++ {
		fmt.Println("%")
		time.Sleep(1 * time.Second)
	}
	timerDone = true
	return nil
}

func Test_Server_GracefulShutdownOngoingWorker(t *testing.T) {
	// This test case explain the minimum/basic workflow of buffalo app.
	r := require.New(t)
	var wg sync.WaitGroup

	// Setup a new buffalo.App with a simple time consuming handler.
	app := New(Options{})
	app.Worker.Register("timer", timerWorker)
	app.Worker.PerformIn(worker.Job{
		Handler: "timer",
	}, 1*time.Second)

	startApp(app, &wg, r) // starts buffalo app routine.

	time.Sleep(1 * time.Second) // make sure just 1 second

	app.cancel()
	time.Sleep(1 * time.Second) // make sure the server started shutdown.

	// This routine is the 2nd client that GETs after Stop it
	// The result should be connection refused even though app is still on.
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := http.Get("http://127.0.0.1:3000")
		r.Contains(err.Error(), "refused")
	}()

	wg.Wait()
	r.Equal(true, timerDone)
}
