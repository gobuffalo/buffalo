package build

import (
	"os/exec"
	"runtime"
	"strings"
)

func buildCmd(opts *Options) (*exec.Cmd, error) {
	if len(opts.GoCommand) == 0 {
		opts.GoCommand = "build"
	}
	buildArgs := []string{opts.GoCommand}

	if len(opts.Mod) != 0 {
		buildArgs = append(buildArgs, "-mod", opts.Mod)
	}

	buildArgs = append(buildArgs, opts.BuildFlags...)

	tf := opts.App.BuildTags(opts.Environment, opts.Tags...)
	if len(tf) > 0 {
		buildArgs = append(buildArgs, "-tags", tf.String())
	}

	if opts.GoCommand == "build" {
		bin := opts.App.Bin
		if runtime.GOOS == "windows" {
			if !strings.HasSuffix(bin, ".exe") {
				bin += ".exe"
			}
			bin = strings.Replace(bin, "/", "\\", -1)
		} else {
			bin = strings.TrimSuffix(bin, ".exe")
		}
		buildArgs = append(buildArgs, "-o", bin)
	}

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

	return exec.Command("go", buildArgs...), nil
}
