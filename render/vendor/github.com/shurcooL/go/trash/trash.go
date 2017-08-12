// +build !darwin

package trash

import "errors"

// MoveTo moves named file or directory to trash.
func MoveTo(name string) error {
	return errors.New("MoveToTrash: not yet implemented on non-darwin")
}
