package gopathutil

import "errors"

// RemoveRepo removes go-gettable repo with no local changes (by moving it into trash).
// importPathPattern must match exactly with the repo root.
// For example, "github.com/user/repo/...".
//
// It's currently not implemented.
func RemoveRepo(importPathPattern string) error {
	return errors.New("not implemented: RemoveRepo needs to be updated to use new dependencies")
}
