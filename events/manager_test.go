package events

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var m, _ = boss.(*manager)

func Test_manager_Listen(t *testing.T) {
	r := require.New(t)

	m.Reset()

	r.Len(m.listeners, 0)

	m.Listen("foo", func(e Event) {})
	r.Len(m.listeners, 1)
	r.NotNil(m.listeners["foo"])

	m.StopListening("foo")
	r.Len(m.listeners, 0)
}

func Test_manager_Emit(t *testing.T) {
	r := require.New(t)

	m.Reset()

	max := 5
	wg := &sync.WaitGroup{}
	wg.Add(max)

	moot := &sync.Mutex{}
	var es []Event
	m.Listen("foo", func(e Event) {
		moot.Lock()
		defer moot.Unlock()
		es = append(es, e)
		wg.Done()
	})

	for i := 0; i < max; i++ {
		err := m.Emit(Event{
			Kind: "FOO",
		})
		r.NoError(err)
	}

	// because wg.Wait can potentially hang here if there's
	// a bug, let's make sure that doesn't happen
	ctx, cf := context.WithTimeout(context.Background(), 2*time.Second)
	var _ = cf // don't want the cf, but lint complains if i don't keep it

	go func() {
		<-ctx.Done()
		if ctx.Err() != nil {
			panic("test ran too long")
		}
	}()
	wg.Wait()
	r.Len(es, max)

	for _, e := range es {
		r.Equal("foo", e.Kind)
	}
}
