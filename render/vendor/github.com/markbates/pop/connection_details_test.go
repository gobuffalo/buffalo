package pop_test

import (
	"testing"

	"github.com/markbates/pop"
	"github.com/stretchr/testify/require"
)

func Test_ConnectionDetails_Finalize(t *testing.T) {
	r := require.New(t)

	cd := &pop.ConnectionDetails{
		URL: "postgres://user:pass@host:port/database",
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(cd.Database, "database")
	r.Equal(cd.Dialect, "postgres")
	r.Equal(cd.Host, "host")
	r.Equal(cd.Password, "pass")
	r.Equal(cd.Port, "port")
	r.Equal(cd.User, "user")
}

func Test_ConnectionDetails_Finalize_UnknownDialect(t *testing.T) {
	r := require.New(t)
	cd := &pop.ConnectionDetails{
		URL: "unknown://user:pass@host:port/database",
	}
	err := cd.Finalize()
	r.Error(err)
}

func Test_ConnectionDetails_Finalize_SQLite(t *testing.T) {
	r := require.New(t)

	cd := &pop.ConnectionDetails{
		URL: "sqlite3:///tmp/foo.db",
	}
	err := cd.Finalize()
	r.NoError(err)

	r.Equal(cd.Database, "/tmp/foo.db")
	r.Equal(cd.Dialect, "sqlite3")
	r.Equal(cd.Host, "")
	r.Equal(cd.Password, "")
	r.Equal(cd.Port, "")
	r.Equal(cd.User, "")
}
