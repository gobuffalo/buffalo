package vcs

import (
	"fmt"
	"os/exec"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packr/v2"
)

// New generator for adding VCS to an application
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if opts.Provider == "none" {
		return g, nil
	}

	box := packr.New("buffalo:genny:vcs", "../vcs/templates")
	s, err := box.FindString("ignore.tmpl")
	if err != nil {
		return g, err
	}

	p := opts.Provider
	n := fmt.Sprintf(".%signore", p)
	g.File(genny.NewFileS(n, s))
	g.Command(exec.Command(p, "init"))

	args := []string{"add", "."}
	if p == "bzr" {
		// Ensure Bazaar is as quiet as Git
		args = append(args, "-q")
	}
	g.Command(exec.Command(p, args...))
	g.Command(exec.Command(p, "commit", "-q", "-m", "Initial Commit"))
	return g, nil
}
