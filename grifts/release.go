package grifts

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Desc("release", "Generates a CHANGELOG and creates a new GitHub release based on what is in the version.go file.")
var _ = grift.Add("release", func(c *grift.Context) error {
	v, err := findVersion()
	if err != nil {
		return err
	}

	err = installBin()
	if err != nil {
		return err
	}

	err = localTest()
	if err != nil {
		return err
	}

	err = dockerTest()
	if err != nil {
		return err
	}

	grift.Run("shoulders", c)

	if err := push(); err != nil {
		return errors.WithStack(err)
	}

	err = tagRelease(v)
	if err != nil {
		return err
	}

	return runReleaser(v)
})

func installBin() error {
	cmd := exec.Command("go", "install", "-v", "./buffalo")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func localTest() error {
	cmd := exec.Command("go", "test", "-tags", "sqlite", "-v", "-race", "./...")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func dockerTest() error {
	cmd := exec.Command("docker", "build", ".")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func tagRelease(v string) error {
	cmd := exec.Command("git", "tag", v)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "push", "origin", "--tags")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runReleaser(v string) error {
	cmd := exec.Command("goreleaser", "--rm-dist")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func push() error {
	cmd := exec.Command("git", "push", "origin", "master")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func findVersion() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	vfile, err := ioutil.ReadFile(filepath.Join(pwd, "buffalo/cmd/version.go"))
	if err != nil {
		return "", err
	}

	//var Version = "v0.4.0"
	re := regexp.MustCompile(`const Version = "(.+)"`)
	matches := re.FindStringSubmatch(string(vfile))
	if len(matches) < 2 {
		return "", errors.New("failed to find the version")
	}
	v := matches[1]
	if strings.Contains(v, "dev") {
		return "", errors.Errorf("version can not be a dev version %s", v)
	}
	if !strings.HasPrefix(v, "v") {
		return "", errors.Errorf("version must match format `v0.0.0`: %s", v)
	}
	return v, nil
}
