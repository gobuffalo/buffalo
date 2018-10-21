package build

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/pkg/errors"
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

	buildArgs = append(buildArgs, "-o", opts.Bin)

	flags := []string{}

	if opts.Static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(opts.LDFlags) > 0 {
		if foundVersion, _ := regexp.MatchString("-X\\s+main.BuildVersion=", opts.LDFlags); foundVersion {
			return nil, errors.New("the ldflag option '-X main.BuildVersion=' is reserved for Buffalo use")
		}
		if foundBuildTime, _ := regexp.MatchString("-X\\s+main.BuildTime=", opts.LDFlags); foundBuildTime {
			return nil, errors.New("the ldflag option '-X main.BuildTime=' is reserved for Buffalo use")
		}
		flags = append(flags, opts.LDFlags)
	}
	if len(flags) > 0 {
		buildArgs = append(buildArgs, "-ldflags", strings.Join(flags, " "))
	}

	return exec.Command(genny.GoBin(), buildArgs...), nil
}
