package pop

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/pkg/errors"
)

func MigrationCreate(path, name, ext string, up, down []byte) error {
	n := time.Now().UTC()
	s := n.Format("20060102150405")

	err := os.MkdirAll(path, 0766)
	if err != nil {
		return errors.Wrapf(err, "couldn't create migrations path %s", path)
	}

	upf := filepath.Join(path, (fmt.Sprintf("%s_%s.up.%s", s, name, ext)))
	err = ioutil.WriteFile(upf, up, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", upf)
	}
	fmt.Printf("> %s\n", upf)

	downf := filepath.Join(path, (fmt.Sprintf("%s_%s.down.%s", s, name, ext)))
	err = ioutil.WriteFile(downf, down, 0666)
	if err != nil {
		return errors.Wrapf(err, "couldn't write up migration %s", downf)
	}

	fmt.Printf("> %s\n", downf)
	return nil
}

// MigrateUp is deprecated, and will be removed in a future version. Use FileMigrator#Up instead.
func (c *Connection) MigrateUp(path string) error {
	warningMsg := "Connection#MigrateUp is deprecated, and will be removed in a future version. Use FileMigrator#Up instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}
	log.Println(warningMsg)

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Up()
}

// MigrateDown is deprecated, and will be removed in a future version. Use FileMigrator#Down instead.
func (c *Connection) MigrateDown(path string, step int) error {
	warningMsg := "Connection#MigrateDown is deprecated, and will be removed in a future version. Use FileMigrator#Down instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}
	log.Println(warningMsg)

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Down(step)
}

// MigrateStatus is deprecated, and will be removed in a future version. Use FileMigrator#Status instead.
func (c *Connection) MigrateStatus(path string) error {
	warningMsg := "Connection#MigrateStatus is deprecated, and will be removed in a future version. Use FileMigrator#Status instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}
	log.Println(warningMsg)

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Status()
}

// MigrateReset is deprecated, and will be removed in a future version. Use FileMigrator#Reset instead.
func (c *Connection) MigrateReset(path string) error {
	warningMsg := "Connection#MigrateReset is deprecated, and will be removed in a future version. Use FileMigrator#Reset instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}
	log.Println(warningMsg)

	mig, err := NewFileMigrator(path, c)
	if err != nil {
		return errors.WithStack(err)
	}
	return mig.Reset()
}
