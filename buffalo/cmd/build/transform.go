package build

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (b *Builder) transform(path string, fn func([]byte, io.Writer) error) error {
	logrus.Debugf("transforming %s", path)
	path = filepath.Join(b.Root, path)
	s, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, ok := b.originals[path]; ok {
		return errors.Errorf("%s was already transformed and can't be altered twice", path)
	}
	b.originals[path] = s
	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	b.cleanups = append(b.cleanups, func() error {
		logrus.Debugf("restoring %s", path)
		f, err := os.Create(path)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		_, err = f.Write(b.originals[path])
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})

	return fn(s, f)
}
