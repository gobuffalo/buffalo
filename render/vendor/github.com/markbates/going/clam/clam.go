package clam

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
)

func RunAndListen(cmd *exec.Cmd, fn func(s string)) error {
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	r, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr.String())
	}

	scanner := bufio.NewScanner(r)
	go func() {
		for scanner.Scan() {
			fn(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr.String())
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr.String())
	}
	return nil
}

func Run(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, stderr.String())
	}

	return out.String(), err
}
