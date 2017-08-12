// +build disable

package gopathutil

import (
	"errors"
	"strings"

	"github.com/kisielk/gotool"
	"github.com/shurcooL/go/gists/gist7480523"
	"github.com/shurcooL/go/trash"
)

// RemoveRepo removes go-gettable repo with no local changes (by moving it into trash).
// importPathPattern must match exactly with the repo root.
// For example, "github.com/user/repo/...".
func RemoveRepo(importPathPattern string) error {
	// TODO: Use an official Go package for `go list` functionality whenever possible.
	importPaths := gotool.ImportPaths([]string{importPathPattern})
	if len(importPaths) == 0 {
		return errors.New("no packages to remove")
	}

	var firstGoPackage *gist7480523.GoPackage
	for i, importPath := range importPaths {
		goPackage := gist7480523.GoPackageFromImportPath(importPath)
		if goPackage == nil {
			return errors.New("Import Path not found: " + importPath)
		}

		if goPackage.Bpkg.Goroot {
			return errors.New("can't remove packages from GOROOT")
		}

		goPackage.UpdateVcs()

		if goPackage.Dir.Repo == nil {
			return errors.New("can't get repo status")
		}

		if i == 0 {
			firstGoPackage = goPackage
		} else if firstGoPackage.Dir.Repo != goPackage.Dir.Repo {
			return errors.New("matched Go Packages span more than 1 repo: " + firstGoPackage.Dir.Repo.Vcs.RootPath() + " != " + goPackage.Dir.Repo.Vcs.RootPath())
		} else if !strings.HasPrefix(goPackage.Bpkg.Dir, firstGoPackage.Dir.Repo.Vcs.RootPath()) { // TODO: This is probably not neccessary...
			return errors.New("Go Package not inside repo: " + goPackage.Bpkg.Dir + " doesn't have prefix " + firstGoPackage.Dir.Repo.Vcs.RootPath())
		}
	}

	if repoImportPathPattern := gist7480523.GetRepoImportPathPattern(firstGoPackage.Dir.Repo.Vcs.RootPath(), firstGoPackage.Bpkg.SrcRoot); repoImportPathPattern != importPathPattern {
		return errors.New("importPathPattern not exact repo root match: " + importPathPattern + " != " + repoImportPathPattern)
	}

	firstGoPackage.UpdateVcsFields()

	cleanStatus := func(goPackage *gist7480523.GoPackage) bool {
		packageStatus := presenter(goPackage)[:4]
		return packageStatus == "    " || packageStatus == "  + " // Updates are okay to ignore.
	}

	if !cleanStatus(firstGoPackage) {
		return errors.New("non-clean status: " + presenter(firstGoPackage))
	}

	err := trash.MoveTo(firstGoPackage.Dir.Repo.Vcs.RootPath())
	return err

	// TODO: Clean up /pkg folder contents, if any, etc.
}

// TODO: Inline and simplify this.
var presenter gist7480523.GoPackageStringer = func(goPackage *gist7480523.GoPackage) string {
	out := ""

	if repo := goPackage.Dir.Repo; repo != nil {
		repoImportPath := gist7480523.GetRepoImportPath(repo.Vcs.RootPath(), goPackage.Bpkg.SrcRoot)

		if repo.VcsLocal.LocalBranch != repo.Vcs.GetDefaultBranch() {
			out += "b"
		} else {
			out += " "
		}
		if repo.VcsLocal.Status != "" {
			out += "*"
		} else {
			out += " "
		}
		if repo.RepoRoot == nil || repo.RepoRoot.Repo != repo.VcsLocal.Remote {
			out += "#"
		} else if repo.VcsLocal.LocalRev != repo.VcsRemote.RemoteRev {
			if repo.VcsRemote.RemoteRev != "" {
				if !repo.VcsRemote.IsContained {
					out += "+"
				} else {
					out += "-"
				}
			} else {
				out += "!"
			}
		} else {
			out += " "
		}
		if repo.VcsLocal.Stash != "" {
			out += "$"
		} else {
			out += " "
		}

		out += " " + repoImportPath + "/..."
	} else {
		out += "????"

		out += " " + goPackage.Bpkg.ImportPath
	}

	return out
}
