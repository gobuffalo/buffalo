package grifts

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/markbates/grift/grift"
)

var _ = grift.Desc("release", "Generates a CHANGELOG and creates a new GitHub release based on what is in the version.go file.")
var _ = grift.Add("release", func(c *grift.Context) error {
	grift.Run("shoulders", c)
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

	err = tagRelease(v)
	if err != nil {
		return err
	}

	err = runChangelogGenerator(v)
	if err != nil {
		return err
	}

	return commitAndPush(v)
})

func installBin() error {
	cmd := exec.Command("go", "install", "-v", "./buffalo")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func localTest() error {
	cmd := exec.Command("go", "test", "-v", "-race", "./...")
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
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return errors.New("GITHUB_TOKEN is not set")
	}

	body := map[string]interface{}{
		"tag_name":   fmt.Sprintf("v%s", v),
		"prerelease": false,
	}

	b, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("https://api.github.com/repos/gobuffalo/buffalo/releases?access_token=%s", token), "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}

	code := res.StatusCode
	if code < 200 || code >= 300 {
		return fmt.Errorf("got a not successful status code from github! %d", code)
	}

	return nil
}

func runChangelogGenerator(v string) error {
	cmd := exec.Command("github_changelog_generator")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func commitAndPush(v string) error {
	cmd := exec.Command("git", "commit", "CHANGELOG.md", "-m", fmt.Sprintf("Updated changelog for release v%s", v))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "push", "origin", "master")
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

	//var Version = "0.4.0"
	re := regexp.MustCompile(`var Version = "(.+)"`)
	matches := re.FindStringSubmatch(string(vfile))
	if len(matches) < 2 {
		return "", errors.New("failed to find the version")
	}
	v := matches[1]
	return v, nil
}
