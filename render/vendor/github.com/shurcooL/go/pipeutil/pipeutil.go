// Package pipeutil provides additional functionality for gopkg.in/pipe.v2 package.
package pipeutil

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"gopkg.in/pipe.v2"
)

// ExecCombinedOutput returns a pipe that runs the named program with the given arguments,
// while forwarding stderr to stdout.
func ExecCombinedOutput(name string, args ...string) pipe.Pipe {
	return func(s *pipe.State) error {
		s.AddTask(&execCombinedOutputTask{name: name, args: args})
		return nil
	}
}

type execCombinedOutputTask struct {
	name string
	args []string

	m      sync.Mutex
	p      *os.Process
	cancel bool
}

func (f *execCombinedOutputTask) Run(s *pipe.State) error {
	f.m.Lock()
	if f.cancel {
		f.m.Unlock()
		return nil
	}
	cmd := exec.Command(f.name, f.args...)
	cmd.Dir = s.Dir
	cmd.Env = s.Env
	cmd.Stdin = s.Stdin
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stdout
	err := cmd.Start()
	f.p = cmd.Process
	f.m.Unlock()
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return &execError{f.name, err}
	}
	return nil
}

func (f *execCombinedOutputTask) Kill() {
	f.m.Lock()
	p := f.p
	f.cancel = true
	f.m.Unlock()
	if p != nil {
		p.Kill()
	}
}

type execError struct {
	name string
	err  error
}

func (e *execError) Error() string {
	return fmt.Sprintf("command %q: %v", e.name, e.err)
}
