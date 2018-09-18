package events

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_manager_Listen(t *testing.T) {
	r := require.New(t)

	boss.Reset()

	r.Len(boss.listeners, 0)

	boss.Listen("foo", func(e Event) {})
	r.Len(boss.listeners, 1)
	r.NotNil(boss.listeners["foo"])

	boss.StopListening("foo")
	r.Len(boss.listeners, 0)
}

func Test_manager_Emit(t *testing.T) {
	r := require.New(t)

	boss.Reset()

	max := 5
	wg := &sync.WaitGroup{}
	wg.Add(max)

	moot := &sync.Mutex{}
	var es []Event
	boss.Listen("foo", func(e Event) {
		moot.Lock()
		defer moot.Unlock()
		es = append(es, e)
		wg.Done()
	})

	for i := 0; i < max; i++ {
		err := boss.Emit(Event{
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
