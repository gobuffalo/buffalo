package build

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

func version() (string, string) {
	_, err := exec.LookPath("git")
	buildTime := fmt.Sprintf("\"%s\"", time.Now().Format(time.RFC3339))
	version := buildTime
	if err == nil {
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		out := &bytes.Buffer{}
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err == nil && out.String() != "" {
			version = strings.TrimSpace(out.String())
		}
	}
	return version, buildTime
}

func (b *Builder) buildBin() error {
	buildArgs := []string{"build", "-i"}
	if b.Debug {
		buildArgs = append(buildArgs, "-v")
	}
	if len(b.Tags) > 0 {
		buildArgs = append(buildArgs, "-tags", strings.Join(b.Tags, " "))
	}

	buildArgs = append(buildArgs, "-o", b.Bin)

	version, buildTime := version()
	flags := []string{
		fmt.Sprintf("-X main.BuildVersion=%s", version),
		fmt.Sprintf("-X main.BuildTime=%s", buildTime),
	}

	if b.Static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(b.LDFlags) > 0 {
		if foundVersion, _ := regexp.MatchString("-X\\s+main.version=", b.LDFlags); foundVersion {
			return errors.New("the ldflag option '-X main.version=' is reserved for Buffalo use")
		}
		if foundBuildTime, _ := regexp.MatchString("-X\\s+main.buildTime=", b.LDFlags); foundBuildTime {
			return errors.New("the ldflag option '-X main.buildTime=' is reserved for Buffalo use")
		}
		flags = append(flags, b.LDFlags)
	}

	buildArgs = append(buildArgs, "-ldflags", strings.Join(flags, " "))

	return b.exec(envy.Get("GO_BIN", "go"), buildArgs...)
}
