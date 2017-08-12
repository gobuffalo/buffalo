package pop

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/pkg/errors"
)

var mrx = regexp.MustCompile("(\\d+)_(.+)\\.(up|down)\\.(sql|fizz)")

func init() {
	MapTableName("schema_migrations", "schema_migration")
	MapTableName("schema_migration", "schema_migration")
}

// NewMigrator returns a new "blank" migrator. It is recommended
// to use something like MigrationBox or FileMigrator. A "blank"
// Migrator should only be used as the basis for a new type of
// migration system.
func NewMigrator(c *Connection) Migrator {
	return Migrator{
		Connection: c,
		Migrations: map[string]Migrations{
			"up":   Migrations{},
			"down": Migrations{},
		},
	}
}

// Migrator forms the basis of all migrations systems.
// It does the actual heavy lifting of running migrations.
// When building a new migration system, you should embed this
// type into your migrator.
type Migrator struct {
	Connection *Connection
	SchemaPath string
	Migrations map[string]Migrations
}

// Up runs pending "up" migrations and applies them to the database.
func (m Migrator) Up() error {
	c := m.Connection
	return m.exec(func() error {
		mfs := m.Migrations["up"]
		sort.Sort(mfs)
		for _, mi := range mfs {
			exists, err := c.Where("version = ?", mi.Version).Exists("schema_migration")
			if err != nil {
				return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
			}
			if exists {
				continue
			}
			err = c.Transaction(func(tx *Connection) error {
				err := mi.Run(tx)
				if err != nil {
					return err
				}
				_, err = tx.Store.Exec(fmt.Sprintf("insert into schema_migration (version) values ('%s')", mi.Version))
				return errors.Wrapf(err, "problem inserting migration version %s", mi.Version)
			})
			if err != nil {
				return errors.WithStack(err)
			}
			fmt.Printf("> %s\n", mi.Name)
		}
		return nil
	})
}

// Down runs pending "down" migrations and rolls back the
// database by the specfied number of steps.
func (m Migrator) Down(step int) error {
	c := m.Connection
	return m.exec(func() error {
		count, err := c.Count("schema_migration")
		if err != nil {
			return errors.Wrap(err, "migration down: unable count existing migration")
		}
		mfs := m.Migrations["down"]
		sort.Sort(sort.Reverse(mfs))
		// skip all runned migration
		if len(mfs) > count {
			mfs = mfs[len(mfs)-count:]
		}
		// run only required steps
		if step > 0 && len(mfs) >= step {
			mfs = mfs[:step]
		}
		for _, mi := range mfs {
			exists, err := c.Where("version = ?", mi.Version).Exists("schema_migration")
			if err != nil || !exists {
				return errors.Wrapf(err, "problem checking for migration version %s", mi.Version)
			}
			err = c.Transaction(func(tx *Connection) error {
				err := mi.Run(tx)
				if err != nil {
					return err
				}
				err = tx.RawQuery("delete from schema_migration where version = ?", mi.Version).Exec()
				return errors.Wrapf(err, "problem deleting migration version %s", mi.Version)
			})
			if err == nil {
				fmt.Printf("< %s\n", mi.Name)
			}
			return err
		}
		return nil
	})
}

// Reset the database by runing the down migrations followed by the up migrations.
func (m Migrator) Reset() error {
	err := m.Down(-1)
	if err != nil {
		return errors.WithStack(err)
	}
	return m.Up()
}

// CreateSchemaMigrations sets up a table to track migrations. This is an idempotent
// operation.
func (m Migrator) CreateSchemaMigrations() error {
	c := m.Connection
	err := c.Open()
	if err != nil {
		return errors.Wrap(err, "could not open connection")
	}
	_, err = c.Store.Exec("select * from schema_migration")
	if err == nil {
		return nil
	}

	return c.Transaction(func(tx *Connection) error {
		smSQL, err := c.Dialect.FizzTranslator().CreateTable(schemaMigrations)
		if err != nil {
			return errors.Wrap(err, "could not build SQL for schema migration table")
		}
		err = tx.RawQuery(smSQL).Exec()
		if err != nil {
			return errors.WithStack(errors.Wrap(err, smSQL))
		}
		return nil
	})
}

// Status prints out the status of applied/pending migrations.
func (m Migrator) Status() error {
	err := m.CreateSchemaMigrations()
	if err != nil {
		return errors.WithStack(err)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Version\tName\tStatus\t")
	for _, mf := range m.Migrations["up"] {
		exists, err := m.Connection.Where("version = ?", mf.Version).Exists("schema_migration")
		if err != nil {
			return errors.Wrapf(err, "problem with migration")
		}
		state := "Pending"
		if exists {
			state = "Applied"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", mf.Version, mf.Name, state)
	}
	return w.Flush()
}

// DumpMigrationSchema will generate a file of the current database schema
// based on the value of Migrator.SchemaPath
func (m Migrator) DumpMigrationSchema() error {
	if m.SchemaPath == "" {
		return nil
	}
	c := m.Connection
	f, err := os.Create(filepath.Join(m.SchemaPath, "schema.sql"))
	if err != nil {
		return errors.WithStack(err)
	}
	err = c.Dialect.DumpSchema(f)
	if err != nil {

		return errors.WithStack(err)
	}
	return nil
}

func (m Migrator) exec(fn func() error) error {
	now := time.Now()
	defer m.DumpMigrationSchema()
	defer printTimer(now)

	err := m.CreateSchemaMigrations()
	if err != nil {
		return errors.Wrap(err, "Migrator: problem creating schema migrations")
	}
	return fn()
}

func printTimer(timerStart time.Time) {
	diff := time.Now().Sub(timerStart).Seconds()
	if diff > 60 {
		fmt.Printf("\n%.4f minutes\n", diff/60)
	} else {
		fmt.Printf("\n%.4f seconds\n", diff)
	}
}
