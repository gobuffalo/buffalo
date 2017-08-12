// +build !nosqlite,!appengine,!appenginevm

package pop

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/markbates/going/defaults"
	. "github.com/markbates/pop/columns"
	"github.com/markbates/pop/fizz"
	"github.com/markbates/pop/fizz/translators"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var _ dialect = &sqlite{}

type sqlite struct {
	gil               *sync.Mutex
	smGil             *sync.Mutex
	ConnectionDetails *ConnectionDetails
}

func (m *sqlite) Details() *ConnectionDetails {
	return m.ConnectionDetails
}

func (m *sqlite) URL() string {
	return m.ConnectionDetails.Database + "?_busy_timeout=5000"
}

func (m *sqlite) MigrationURL() string {
	return m.ConnectionDetails.URL
}

func (m *sqlite) Create(s store, model *Model, cols Columns) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericCreate(s, model, cols), "sqlite create")
	})
}

func (m *sqlite) Update(s store, model *Model, cols Columns) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericUpdate(s, model, cols), "sqlite update")
	})
}

func (m *sqlite) Destroy(s store, model *Model) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericDestroy(s, model), "sqlite destroy")
	})
}

func (m *sqlite) SelectOne(s store, model *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericSelectOne(s, model, query), "sqlite select one")
	})
}

func (m *sqlite) SelectMany(s store, models *Model, query Query) error {
	return m.locker(m.smGil, func() error {
		return errors.Wrap(genericSelectMany(s, models, query), "sqlite select many")
	})
}

func (m *sqlite) Lock(fn func() error) error {
	return m.locker(m.gil, fn)
}

func (m *sqlite) locker(l *sync.Mutex, fn func() error) error {
	if defaults.String(m.Details().Options["lock"], "true") == "true" {
		defer l.Unlock()
		l.Lock()
	}
	err := fn()
	attempts := 0
	for err != nil && err.Error() == "database is locked" && attempts <= m.Details().RetryLimit() {
		time.Sleep(m.Details().RetrySleep())
		err = fn()
		attempts++
	}
	return err
}

func (m *sqlite) CreateDB() error {
	d := filepath.Dir(m.ConnectionDetails.Database)
	err := os.MkdirAll(d, 0766)
	if err != nil {
		return errors.Wrapf(err, "could not create SQLite database %s", m.ConnectionDetails.Database)
	}
	fmt.Printf("created database %s\n", m.ConnectionDetails.Database)
	return nil
}

func (m *sqlite) DropDB() error {
	err := os.Remove(m.ConnectionDetails.Database)
	if err != nil {
		return errors.Wrapf(err, "could not drop SQLite database %s", m.ConnectionDetails.Database)
	}
	fmt.Printf("dropped database %s\n", m.ConnectionDetails.Database)
	return nil
}

func (m *sqlite) TranslateSQL(sql string) string {
	return sql
}

func (m *sqlite) FizzTranslator() fizz.Translator {
	return translators.NewSQLite(m.Details().Database)
}

func (m *sqlite) DumpSchema(w io.Writer) error {
	cmd := exec.Command("sqlite3", m.Details().Database, ".schema")
	Log(strings.Join(cmd.Args, " "))
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("dumped schema for %s\n", m.Details().Database)
	return nil
}

func (m *sqlite) LoadSchema(r io.Reader) error {
	cmd := exec.Command("sqlite3", m.ConnectionDetails.Database)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		io.Copy(in, r)
	}()
	Log(strings.Join(cmd.Args, " "))
	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	fmt.Printf("loaded schema for %s\n", m.Details().Database)
	return nil
}

func (m *sqlite) TruncateAll(tx *Connection) error {
	const tableNames = `SELECT name FROM sqlite_master WHERE type = "table"`
	names := []struct {
		Name string `db:"name"`
	}{}

	err := tx.RawQuery(tableNames).All(&names)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return nil
	}
	stmts := []string{}
	for _, n := range names {
		stmts = append(stmts, fmt.Sprintf("DELETE FROM %s", n.Name))
	}
	return tx.RawQuery(strings.Join(stmts, "; ")).Exec()
}

func newSQLite(deets *ConnectionDetails) (dialect, error) {
	deets.URL = fmt.Sprintf("sqlite3://%s", deets.Database)
	cd := &sqlite{
		gil:               &sync.Mutex{},
		smGil:             &sync.Mutex{},
		ConnectionDetails: deets,
	}

	return cd, nil
}
