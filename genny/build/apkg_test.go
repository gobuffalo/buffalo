package build

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_copyDatabase(t *testing.T) {
	r := require.New(t)

	run := cokeRunner()

	r.NoError(copyDatabase(run))

	f, err := run.Disk.Find("a/database.go")
	r.NoError(err)

	r.Contains(f.String(), "var DB_CONFIG = `development")
}

func Test_copyDatabase_Off(t *testing.T) {
	r := require.New(t)

	run := cokeRunner()
	run.Disk.Remove("database.yml")

	r.NoError(copyDatabase(run))

	f, err := run.Disk.Find("a/database.go")
	r.NoError(err)

	r.NotContains(f.String(), "var DB_CONFIG = `development")
	r.Contains(f.String(), "var DB_CONFIG = ``")
}
