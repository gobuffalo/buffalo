package build

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
)

func buildCmd(opts *Options) (*exec.Cmd, error) {
	buildArgs := []string{"build"}

	if !gomods.On() {
		buildArgs = append(buildArgs, "-i")
	}

	buildArgs = append(buildArgs, opts.BuildFlags...)

	tf := opts.App.BuildTags(opts.Environment, opts.Tags...)
	if len(tf) > 0 {
		buildArgs = append(buildArgs, "-tags", tf.String())
	}

	bin := opts.App.Bin
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") {
			bin += ".exe"
		}
		bin = strings.Replace(bin, "/", "\\", -1)
	}
	buildArgs = append(buildArgs, "-o", bin)

	flags := []string{}

	if opts.Static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(opts.LDFlags) > 0 {
		flags = append(flags, opts.LDFlags)
	}
	if len(flags) > 0 {
		buildArgs = append(buildArgs, "-ldflags", strings.Join(flags, " "))
	}

	return exec.Command(genny.GoBin(), buildArgs...), nil
}
