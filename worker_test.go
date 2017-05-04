package buffalo

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_defaultWorker_Perform(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := defaultWorker{}
	w.Perform(func() error {
		hit = true
		wg.Done()
		return nil
	})
	wg.Wait()
	r.True(hit)
}

func Test_defaultWorker_PerformAt(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := defaultWorker{}
	w.PerformAt(func() error {
		hit = true
		wg.Done()
		return nil
	}, time.Now().Add(10*time.Millisecond))
	wg.Wait()
	r.True(hit)
}

func Test_defaultWorker_PerformIn(t *testing.T) {
	r := require.New(t)

	var hit bool
	wg := &sync.WaitGroup{}
	wg.Add(1)
	w := defaultWorker{}
	w.PerformIn(func() error {
		hit = true
		wg.Done()
		return nil
	}, time.Duration(10*time.Millisecond))
	wg.Wait()
	r.True(hit)
}
