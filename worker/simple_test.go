package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func sampleHandler(Args) error {
	return nil
}

func Test_Simple_RegisterEmpty(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	err := w.Register("", sampleHandler)
	r.Error(err)
}

func Test_Simple_RegisterNil(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	err := w.Register("sample", nil)
	r.Error(err)
}

func Test_Simple_RegisterEmptyNil(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	err := w.Register("", nil)
	r.Error(err)
}

func Test_Simple_RegisterExisting(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	err := w.Register("sample", sampleHandler)
	r.NoError(err)

	err = w.Register("sample", sampleHandler)
	r.Error(err)
}

func Test_Simple_StartStop(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	ctx := context.Background()
	err := w.Start(ctx)
	r.NoError(err)
	r.NotNil(w.ctx)
	r.Nil(w.ctx.Err())

	err = w.Stop()
	r.NoError(err)
	r.NotNil(w.ctx)
	r.NotNil(w.ctx.Err())
}

func Test_Simple_Perform(t *testing.T) {
	r := require.New(t)

	var hit bool
	w := NewSimple()
	r.NoError(w.Start(context.Background()))

	w.Register("x", func(Args) error {
		hit = true
		return nil
	})
	w.Perform(Job{
		Handler: "x",
	})

	// the worker should guarantee the job is finished before the worker stopped
	r.NoError(w.Stop())
	r.True(hit)
}

func Test_Simple_PerformBroken(t *testing.T) {
	r := require.New(t)

	var hit bool
	w := NewSimple()
	r.NoError(w.Start(context.Background()))

	w.Register("x", func(Args) error {
		hit = true

		//Index out of bounds on purpose
		println([]string{}[0])

		return nil
	})
	w.Perform(Job{
		Handler: "x",
	})

	r.NoError(w.Stop())
	r.True(hit)
}

func Test_Simple_PerformWithEmptyJob(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	r.NoError(w.Start(context.Background()))
	defer w.Stop()

	err := w.Perform(Job{})
	r.Error(err)
}

func Test_Simple_PerformWithUnknownJob(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	r.NoError(w.Start(context.Background()))
	defer w.Stop()

	err := w.Perform(Job{Handler: "unknown"})
	r.Error(err)
}

/* TODO(sio4): #road-to-v1 - define the purpose of Start clearly
   consider to make Perform to work only when the worker is started.
func Test_Simple_PerformBeforeStart(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	r.NoError(w.Register("sample", sampleHandler))

	err := w.Perform(Job{Handler: "sample"})
	r.Error(err)
}
*/

func Test_Simple_PerformAfterStop(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	r.NoError(w.Register("sample", sampleHandler))
	r.NoError(w.Start(context.Background()))
	r.NoError(w.Stop())

	err := w.Perform(Job{Handler: "sample"})
	r.Error(err)
}

func Test_Simple_PerformAt(t *testing.T) {
	r := require.New(t)

	var hit bool
	w := NewSimple()
	r.NoError(w.Start(context.Background()))

	w.Register("x", func(Args) error {
		hit = true
		return nil
	})
	w.PerformAt(Job{
		Handler: "x",
	}, time.Now().Add(5*time.Millisecond))

	time.Sleep(10 * time.Millisecond)
	r.True(hit)

	r.NoError(w.Stop())
}

func Test_Simple_PerformIn(t *testing.T) {
	r := require.New(t)

	var hit bool
	w := NewSimple()
	r.NoError(w.Start(context.Background()))

	w.Register("x", func(Args) error {
		hit = true
		return nil
	})
	w.PerformIn(Job{
		Handler: "x",
	}, 5*time.Millisecond)

	time.Sleep(10 * time.Millisecond)
	r.True(hit)

	r.NoError(w.Stop())
}

/* TODO(sio4): #road-to-v1 - define the purpose of Start clearly
   consider to make Perform to work only when the worker is started.
func Test_Simple_PerformInBeforeStart(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	r.NoError(w.Register("sample", sampleHandler))

	err := w.PerformIn(Job{Handler: "sample"}, 5*time.Millisecond)
	r.Error(err)
}
*/

func Test_Simple_PerformInAfterStop(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	r.NoError(w.Register("sample", sampleHandler))
	r.NoError(w.Start(context.Background()))
	r.NoError(w.Stop())

	err := w.PerformIn(Job{Handler: "sample"}, 5*time.Millisecond)
	r.Error(err)
}

/* TODO(sio4): #road-to-v1 - it should be guaranteed to be performed
   consider to make PerformIn to guarantee the job execution
func Test_Simple_PerformInFollowedByStop(t *testing.T) {
	r := require.New(t)

	var hit bool
	w := NewSimple()
	r.NoError(w.Start(context.Background()))

	w.Register("x", func(Args) error {
		hit = true
		return nil
	})
	err := w.PerformIn(Job{
		Handler: "x",
	}, 5*time.Millisecond)
	r.NoError(err)

	// stop the worker immediately after PerformIn
	r.NoError(w.Stop())

	r.True(hit)
}
*/
