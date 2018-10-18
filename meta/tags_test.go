package meta

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_BuildTags(t *testing.T) {
	defer os.Remove("database.yml")
	app := New(".")
	t.Run("without sqlite", func(st *testing.T) {
		r := require.New(st)
		tags := app.BuildTags("dev")
		r.Len(tags, 1)
		r.Equal("dev", tags[0])
		r.Equal(`dev`, tags.String())
	})
	t.Run("with sqlite", func(st *testing.T) {
		r := require.New(st)
		app.WithSQLite = true

		tags := app.BuildTags("dev")
		r.Len(tags, 2)
		r.Equal("dev", tags[0])
		r.Equal("sqlite", tags[1])
		r.Equal(`dev sqlite`, tags.String())
	})
}
