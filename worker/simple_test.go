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
	w := Simple{}
	w.Perform(Job{
		Handler: func(...interface{}) error {
			hit = true
			wg.Done()
			return nil
		},
	})
	wg.Wait()
	r.True(hit)
}

func Test_Simple_PerformAt(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := Simple{}
	w.PerformAt(Job{
		Handler: func(...interface{}) error {
			hit = true
			wg.Done()
			return nil
		},
	}, time.Now().Add(5*time.Millisecond))
	wg.Wait()
	r.True(hit)
}

func Test_Simple_PerformIn(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := Simple{}
	w.PerformIn(Job{
		Handler: func(args ...interface{}) error {
			hit = true
			wg.Done()
			return nil
		},
	}, 5*time.Millisecond)
	wg.Wait()
	r.True(hit)
}

func Test_Simple_NoHandler(t *testing.T) {
	r := require.New(t)

	w := Simple{}
	err := w.Perform(Job{})
	r.Error(err)
}
