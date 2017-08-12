// +build nosqlite appengine appenginevm

package pop

import "errors"

func newSQLite(deets *ConnectionDetails) (dialect, error) {
	return nil, errors.New("sqlite3 was not compiled into the binary")
}
