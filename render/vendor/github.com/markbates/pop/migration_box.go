package pop

import (
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// MigrationBox is a wrapper around packr.Box and Migrator.
// This will allow you to run migrations from a packed box
// inside of a compiled binary.
type MigrationBox struct {
	Migrator
	Box packr.Box
}

// NewMigrationBox from a packr.Box and a Connection.
func NewMigrationBox(box packr.Box, c *Connection) (MigrationBox, error) {
	fm := MigrationBox{
		Migrator: NewMigrator(c),
		Box:      box,
	}

	err := fm.findMigrations()
	if err != nil {
		return fm, errors.WithStack(err)
	}

	return fm, nil
}

func (fm *MigrationBox) findMigrations() error {
	return fm.Box.Walk(func(p string, f packr.File) error {
		info, err := f.FileInfo()
		if err != nil {
			return errors.WithStack(err)
		}
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
		return nil
	})
}
