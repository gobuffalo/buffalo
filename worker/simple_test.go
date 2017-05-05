package worker

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Simple_Perform(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := NewSimple()
	w.Register("x", func(Args) error {
		hit = true
		wg.Done()
		return nil
	})
	w.Perform(Job{
		Handler: "x",
	})
	wg.Wait()
	r.True(hit)
}

func Test_Simple_PerformAt(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := NewSimple()
	w.Register("x", func(Args) error {
		hit = true
		wg.Done()
		return nil
	})
	w.PerformAt(Job{
		Handler: "x",
	}, time.Now().Add(5*time.Millisecond))
	wg.Wait()
	r.True(hit)
}

func Test_Simple_PerformIn(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := NewSimple()
	w.Register("x", func(Args) error {
		hit = true
		wg.Done()
		return nil
	})
	w.PerformIn(Job{
		Handler: "x",
	}, 5*time.Millisecond)
	wg.Wait()
	r.True(hit)
}

func Test_Simple_NoHandler(t *testing.T) {
	r := require.New(t)

	w := NewSimple()
	err := w.Perform(Job{})
	r.Error(err)
}
