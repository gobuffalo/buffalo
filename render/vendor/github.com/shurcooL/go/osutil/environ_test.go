package osutil_test

import (
	"os"
	"os/exec"

	"github.com/shurcooL/go/osutil"
)

func ExampleEnviron() {
	cmd := exec.Command("example")
	env := osutil.Environ(os.Environ())
	env.Set("USER", "gopher")
	env.Set("HOME", "/usr/gopher")
	env.Unset("TMPDIR")
	cmd.Env = env
}
