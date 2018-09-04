package grifts

import (
	"strings"

	"github.com/gobuffalo/buffalo/grifts/internal/release"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
)

var _ = grift.Desc("release", "Generates a CHANGELOG and creates a new GitHub release based on what is in the version.go file.")
var _ = grift.Add("release", func(c *grift.Context) error {

	v, err := release.FindVersion("runtime/version.go")
	if err != nil {
		return err
	}

	rr, err := release.New(v)
	if err != nil {
		return err
	}

	rr.Add(release.InBranch("master", func() error {
		m, err := release.New(v)
		if err != nil {
			return errors.WithStack(err)
		}
		m.Add(release.Command("make", "install"))
		m.Add(release.Command("make", "ci-test"))
		m.Add(release.Runner{
			Name: "grift shoulders",
			Fn: func() error {
				return grift.Run("shoulders", c)
			},
		})
		p, err := release.PackAndCommit()
		if err != nil {
			if !strings.Contains(err.Error(), "nothing to commit, working tree clean") {
				return errors.WithStack(err)
			}
		}
		m.Add(p)

		tr, err := release.TagRelease("master", v)
		if err != nil {
			return errors.WithStack(err)
		}
		m.Add(tr)
		m.Add(release.Command("bash", "'goreleaser --rm-dist'"))

		m.Add(release.Command("git", "branch", "-D", "development"))
		m.Add(release.Command("git", "branch", "development"))
		return m.Run()
	}))

	rr.Add(release.InBranch("development", func() error {
		m, err := release.New("development")
		if err != nil {
			return errors.WithStack(err)
		}
		p, err := release.UnpackAndCommit()
		if err != nil {
			return errors.WithStack(err)
		}
		m.Add(p)
		return m.Run()
	}))
	return rr.Run()
})
