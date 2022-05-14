package buffalo

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/gobuffalo/buffalo/plugins"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/events"
)

var loadPlugins sync.Once

// LoadPlugins will add listeners for any plugins that support "events"
func LoadPlugins() error {
	var errResult error
	loadPlugins.Do(func() {
		// don't send plugins events during testing
		if envy.Get("GO_ENV", "development") == "test" {
			return
		}
		plugs, err := plugins.Available()
		if err != nil {
			errResult = err
			return
		}
		for _, cmds := range plugs {
			for _, c := range cmds {
				if c.BuffaloCommand != "events" {
					continue
				}
				err := func(c plugins.Command) error {
					n := fmt.Sprintf("[PLUGIN] %s %s", c.Binary, c.Name)
					fn := func(e events.Event) {
						b, err := json.Marshal(e)
						if err != nil {
							fmt.Println("error trying to marshal event", e, err)
							return
						}
						cmd := exec.Command(c.Binary, c.UseCommand, string(b))
						cmd.Stderr = os.Stderr
						cmd.Stdout = os.Stdout
						cmd.Stdin = os.Stdin
						if err := cmd.Run(); err != nil {
							fmt.Println("error trying to send event", strings.Join(cmd.Args, " "), err)
						}
					}
					_, err := events.NamedListen(n, events.Filter(c.ListenFor, fn))
					return err
				}(c)
				if err != nil {
					errResult = err
					return
				}
			}
		}
	})
	return errResult
}
