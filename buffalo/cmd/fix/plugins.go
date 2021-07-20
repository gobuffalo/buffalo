package fix

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	cmdPlugins "github.com/gobuffalo/buffalo/buffalo/cmd/plugins"
	"github.com/gobuffalo/buffalo/genny/plugins/install"
	"github.com/gobuffalo/buffalo/internal/takeon/github.com/markbates/errx"
	"github.com/gobuffalo/buffalo/plugins"
	"github.com/gobuffalo/buffalo/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
)

//Plugins fixes the plugin configuration of the project by
//manipulating the plugins .toml file.
type Plugins struct{}

//CleanCache cleans the plugins cache folder by removing it
func (pf Plugins) CleanCache(r *Runner) error {
	fmt.Println("~~~ Cleaning plugins cache ~~~")
	os.RemoveAll(plugins.CachePath)
	return nil
}

//Reinstall installs latest versions of the plugins
func (pf Plugins) Reinstall(r *Runner) error {
	plugs, err := plugdeps.List(r.App)
	if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
		return err
	}

	run := genny.WetRunner(context.Background())
	gg, err := install.New(&install.Options{
		App:     r.App,
		Plugins: plugs.List(),
	})
	if err != nil {
		return err
	}

	run.WithGroup(gg)

	fmt.Println("~~~ Reinstalling plugins ~~~")
	return run.Run()
}

//RemoveOld removes old and deprecated plugins
func (pf Plugins) RemoveOld(r *Runner) error {
	fmt.Println("~~~ Removing old plugins ~~~")

	run := genny.WetRunner(context.Background())
	app := meta.New(".")
	plugs, err := plugdeps.List(app)
	if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
		return err
	}

	a := strings.TrimSpace("github.com/gobuffalo/buffalo-pop")
	bin := path.Base(a)
	plugs.Remove(plugdeps.Plugin{
		Binary: bin,
		GoGet:  a,
	})

	fmt.Println("~~~ Removing github.com/gobuffalo/buffalo-pop plugin ~~~")

	run.WithRun(cmdPlugins.NewEncodePluginsRunner(app, plugs))

	return run.Run()
}
