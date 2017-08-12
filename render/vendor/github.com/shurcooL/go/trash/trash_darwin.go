package trash

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// MoveTo moves named file or directory to trash.
func MoveTo(name string) error {
	name = filepath.Clean(name)
	home := os.Getenv("HOME")
	dir, file := filepath.Split(name)
	target := filepath.Join(home, ".Trash", file)

	// TODO: If target name exists in Trash, come up with a unique one (perhaps append a timestamp) instead of overwriting.
	// TODO: Support OS X "Put Back". Figure out how it's done and do it.

	err := os.Rename(name, target)
	if err != nil {
		return err
	}

	// If directory became empty, remove it (recursively up).
	for {
		// Ensure it's an empty directory.
		if dirEntries, err := ioutil.ReadDir(dir); err != nil || len(dirEntries) != 0 {
			break
		}

		// Remove directory if it's (now) empty.
		err := os.Remove(dir)
		if err != nil {
			break
		}

		dir, _ = filepath.Split(dir)
	}

	return nil
}
