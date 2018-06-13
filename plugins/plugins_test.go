package plugins

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/gobuffalo/envy"
)

func TestAskBin_respectsTimeout(t *testing.T) {
	from, err := envy.MustGet("BUFFALO_PLUGIN_PATH")
	if err != nil {
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
		t.Error(err)
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
		t.Log("did not time-out quickly enough")
		t.Fail()
	case <-done:
		t.Log("timed-out successfully")
	}
}
