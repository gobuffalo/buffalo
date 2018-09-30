package build

import (
	"regexp"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/pkg/errors"
)

func (b *Builder) buildBin() error {
	buildArgs := []string{"build"}
	if b.Debug {
		buildArgs = append(buildArgs, "-v")
	}

	if !gomods.On() {
		buildArgs = append(buildArgs, "-i")
	}

	tf := b.App.BuildTags(b.Environment, b.Tags...)
	buildArgs = append(buildArgs, "-tags", tf.String())

	buildArgs = append(buildArgs, "-o", b.Bin)

	flags := []string{}

	if b.Static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(b.LDFlags) > 0 {
		if foundVersion, _ := regexp.MatchString("-X\\s+main.BuildVersion=", b.LDFlags); foundVersion {
			return errors.New("the ldflag option '-X main.BuildVersion=' is reserved for Buffalo use")
		}
		if foundBuildTime, _ := regexp.MatchString("-X\\s+main.BuildTime=", b.LDFlags); foundBuildTime {
			return errors.New("the ldflag option '-X main.BuildTime=' is reserved for Buffalo use")
		}
		flags = append(flags, b.LDFlags)
	}
	if len(flags) > 0 {
		buildArgs = append(buildArgs, "-ldflags", strings.Join(flags, " "))
	}

	return b.exec(envy.Get("GO_BIN", "go"), buildArgs...)
}
