package pop

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Store is an interface that must be implemented in order for Pop
// to be able to use the value as a way of talking to a datastore.
type store interface {
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	NamedExec(string, interface{}) (sql.Result, error)
	Exec(string, ...interface{}) (sql.Result, error)
	PrepareNamed(string) (*sqlx.NamedStmt, error)
	Transaction() (*tX, error)
	Rollback() error
	Commit() error
	Close() error
}
