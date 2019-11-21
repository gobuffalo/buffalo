package plugins

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAskBin_respectsTimeout(t *testing.T) {
	r := require.New(t)

	from := os.Getenv(BUFFALO_PLUGIN_PATH)
	if len(from) == 0 {
		t.Skipf("BUFFALO_PLUGIN_PATH not set.")
		return
	}

	if fileEntries, err := ioutil.ReadDir(from); err == nil {
		found := false
		for _, e := range fileEntries {
			if strings.HasPrefix(e.Name(), "buffalo-") {
				from = e.Name()
				found = true
				break
			}
		}
		if !found {
			t.Skipf("no plugins found")
			return
		}
	} else {
		r.Error(err, "plugin path not able to be read")
		return
	}

	const tooShort = time.Millisecond
	impossible, cancel := context.WithTimeout(context.Background(), tooShort)
	defer cancel()

	done := make(chan struct{})
	go func() {
		askBin(impossible, from)
		close(done)
	}()

	select {
	case <-time.After(tooShort + 80*time.Millisecond):
		r.Fail("did not time-out quickly enough")
	case <-done:
		t.Log("timed-out successfully")
	}

	if _, ok := findInCache(from); ok {
		r.Fail("expected plugin not to be added to cache on failure, but it was in cache")
	}
}
