package buffalo

import (
	"context"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gobuffalo/events"
	"github.com/stretchr/testify/require"
)

func Test_Logger_Emits_Events(t *testing.T) {
	r := require.New(t)

	l := NewLogger("debug")

	wg := &sync.WaitGroup{}
	var es []string
	df, err := events.Listen(func(e events.Event) {
		if strings.HasPrefix(e.Kind, "buffalo:log") {
			es = append(es, e.Kind)
			wg.Done()
		}
	})
	r.NoError(err)
	defer df()

	max := 4
	wg.Add(max)
	l.Debug("debugging")
	l.Error("erroring")
	l.Info("infoing")
	l.Warn("warning")

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
	sort.Strings(es)
	r.Equal([]string{"buffalo:log:debug", "buffalo:log:error", "buffalo:log:info", "buffalo:log:warning"}, es)
}
