package pop

import "github.com/pkg/errors"

type Migration struct {
	// Path to the migration (./migrations/123_create_widgets.up.sql)
	Path string
	// Version of the migration (123)
	Version string
	// Name of the migration (create_widgets)
	Name string
	// Direction of the migration (up)
	Direction string
	// Type of migration (sql)
	Type string
	// Runner function to run/execute the migration
	Runner func(Migration, *Connection) error
}

// Run the migration. Returns an error if there is
// no mf.Runner defined.
func (mf Migration) Run(c *Connection) error {
	if mf.Runner == nil {
		return errors.Errorf("no runner defined for %s", mf.Path)
	}
	return mf.Runner(mf, c)
}

type Migrations []Migration

func (mfs Migrations) Len() int {
	return len(mfs)
}

func (mfs Migrations) Less(i, j int) bool {
	return mfs[i].Version < mfs[j].Version
}

func (mfs Migrations) Swap(i, j int) {
	mfs[i], mfs[j] = mfs[j], mfs[i]
}
