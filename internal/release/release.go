package release

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Runner struct {
	Name string
	Fn   func() error
}

type Manager struct {
	Version string
	Runners []Runner
}

func (m *Manager) Add(r Runner) {
	m.Runners = append(m.Runners, r)
}

func (m Manager) Run() error {
	for _, r := range m.Runners {
		fmt.Println(r.Name)
		if err := r.Fn(); err != nil {
			return errors.Wrapf(err, "releaser runner %s failed", r.Name)
		}
	}
	return nil
}

func New(version string) (*Manager, error) {
	version = strings.TrimSpace(version)
	if version == "" {
		return nil, errors.New("version can not be empty")
	}
	return &Manager{
		Version: version,
		Runners: []Runner{},
	}, nil
}

func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "--short", "HEAD")
	b, err := cmd.CombinedOutput()
	return string(b), err
}

func InBranch(name string, fn func() error) Runner {
	f := func() error {
		cur, err := CurrentBranch()
		if err != nil {
			return errors.WithStack(err)
		}
		defer func() {
			if err := Command("git", "checkout", cur).Fn(); err != nil {
				log.Fatal(err)
			}
		}()
		if err := Command("git", "checkout", name).Fn(); err != nil {
			return errors.WithStack(err)
		}
		return fn()
	}
	return Runner{
		Name: "git branch: " + name,
		Fn:   f,
	}
}

func TagRelease(branch string, version string) (Runner, error) {
	r := Runner{
		Name: fmt.Sprintf("tag %s as %s", branch, version),
	}
	m, err := New(version)
	if err != nil {
		return r, errors.WithStack(err)
	}
	m.Add(Command("git", "tag", version))
	m.Add(Command("git", "push", "origin", branch))
	m.Add(Command("git", "push", "origin", "--tags"))
	r.Fn = m.Run
	return r, nil
}

func PackAndCommit() (Runner, error) {
	r := Runner{Name: "packr"}
	m, err := New("n/a")
	if err != nil {
		return r, errors.WithStack(err)
	}
	m.Add(Command("packr"))
	m.Add(Command("git", "add", "**/*-packr.go"))
	m.Add(Command("git", "commit", "**/*-packr.go", "-m", "committed packr files"))
	r.Fn = m.Run
	return r, nil
}

func UnpackAndCommit() (Runner, error) {
	r := Runner{Name: "packr"}
	m, err := New("n/a")
	if err != nil {
		return r, errors.WithStack(err)
	}
	m.Add(Command("packr", "clean"))
	m.Add(Command("git", "add", "**/*-packr.go"))
	m.Add(Command("git", "commit", "**/*-packr.go", "-m", "deleted packr files"))
	r.Fn = m.Run
	return r, nil
}

func Command(name string, args ...string) Runner {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	fn := func() error {
		return cmd.Run()
	}
	r := Runner{
		Name: strings.Join(cmd.Args, " "),
		Fn:   fn,
	}
	return r
}

func FindVersion(path string) (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	vfile, err := ioutil.ReadFile(filepath.Join(pwd, path))
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
