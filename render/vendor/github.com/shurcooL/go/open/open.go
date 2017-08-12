// Package open offers ability to open files or URLs as if user double-clicked it in their OS.
//
// Deprecated: Use github.com/shurcooL/go/browser package if you need to open URLs.
// It respects the BROWSER environment variable.
package open

import (
	"log"
	"os/exec"
	"runtime"
)

// Open opens a file (or a directory or url), just as if the user had double-clicked the file's icon.
// It uses the default application, as determined by the OS.
func Open(path string) {
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open", path}
	case "windows":
		args = []string{"cmd", "/c", "start", path}
	default:
		args = []string{"xdg-open", path}
	}
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Println("open.Open:", err)
	}
}
