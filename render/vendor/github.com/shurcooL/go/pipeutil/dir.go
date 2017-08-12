package pipeutil

import "gopkg.in/pipe.v2"

// OutputDir is identical to pipe.Output, except it sets the starting dir.
func OutputDir(p pipe.Pipe, dir string) ([]byte, error) {
	outb := &pipe.OutputBuffer{}
	s := pipe.NewState(outb, nil)
	s.Dir = dir
	err := p(s)
	if err == nil {
		err = s.RunTasks()
	}
	return outb.Bytes(), err
}

// DividedOutputDir is identical to pipe.DividedOutput, except it sets the starting dir.
func DividedOutputDir(p pipe.Pipe, dir string) (stdout []byte, stderr []byte, err error) {
	outb := &pipe.OutputBuffer{}
	errb := &pipe.OutputBuffer{}
	s := pipe.NewState(outb, errb)
	s.Dir = dir
	err = p(s)
	if err == nil {
		err = s.RunTasks()
	}
	return outb.Bytes(), errb.Bytes(), err
}
