package pop

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/markbates/pop/fizz"
	"github.com/pkg/errors"
)

// FileMigrator is a migrator for SQL and Fizz
// files on disk at a specified path.
type FileMigrator struct {
	Migrator
	Path string
}

// NewFileMigrator for a path and a Connection
func NewFileMigrator(path string, c *Connection) (FileMigrator, error) {
	fm := FileMigrator{
		Migrator: NewMigrator(c),
		Path:     path,
	}

	err := fm.findMigrations()
	if err != nil {
		return fm, errors.WithStack(err)
	}

	return fm, nil
}

func (fm *FileMigrator) findMigrations() error {
	dir := filepath.Base(fm.Path)
	if _, err := os.Stat(dir); err != nil {
		// directory doesn't exist
		return nil
	}
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			matches := mrx.FindAllStringSubmatch(info.Name(), -1)
			if matches == nil || len(matches) == 0 {
				return nil
			}
			m := matches[0]
			mf := Migration{
				Path:      p,
				Version:   m[1],
				Name:      m[2],
				Direction: m[3],
				Type:      m[4],
				Runner: func(mf Migration, tx *Connection) error {
					f, err := os.Open(p)
					if err != nil {
						return errors.WithStack(err)
					}
					content, err := migrationContent(mf, tx, f)
					if err != nil {
						return errors.Wrapf(err, "error processing %s", mf.Path)
					}

					if content == "" {
						return nil
					}

					err = tx.RawQuery(content).Exec()
					if err != nil {
						return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
					}
					return nil
				},
			}
			fm.Migrations[mf.Direction] = append(fm.Migrations[mf.Direction], mf)
		}
		return nil
	})
	return nil
}

func migrationContent(mf Migration, c *Connection, r io.Reader) (string, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", nil
	}

	content := string(b)

	t := template.Must(template.New("sql").Parse(content))
	var bb bytes.Buffer
	err = t.Execute(&bb, c.Dialect.Details())
	if err != nil {
		return "", errors.Wrapf(err, "could not execute migration template %s", mf.Path)
	}
	content = bb.String()

	if mf.Type == "fizz" {
		content, err = fizz.AString(content, c.Dialect.FizzTranslator())
		if err != nil {
			return "", errors.Wrapf(err, "could not fizz the migration %s", mf.Path)
		}
	}
	return content, nil
}
