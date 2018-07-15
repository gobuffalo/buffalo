package meta

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_BuildTags(t *testing.T) {
	defer os.Remove("database.yml")
	app := New(".")
	t.Run("without database.yml", func(st *testing.T) {
		r := require.New(st)
		tags := app.BuildTags("dev")
		r.Len(tags, 1)
		r.Equal("dev", tags[0])
		r.Equal(`dev`, tags.String())
	})
	t.Run("with database.yml", func(st *testing.T) {
		t.Run("with sqlite", func(st *testing.T) {
			r := require.New(st)
			f, err := os.Create("database.yml")
			r.NoError(err)
			_, err = f.WriteString("sqlite")
			r.NoError(err)

			tags := app.BuildTags("dev")
			r.Len(tags, 2)
			r.Equal("dev", tags[0])
			r.Equal("sqlite", tags[1])
			r.Equal(`dev sqlite`, tags.String())
		})
		t.Run("without sqlite", func(st *testing.T) {
			r := require.New(st)
			f, err := os.Create("database.yml")
			r.NoError(err)
			_, err = f.WriteString("mysql")
			r.NoError(err)

			tags := app.BuildTags("dev")
			r.Len(tags, 1)
			r.Equal("dev", tags[0])
			r.Equal(`dev`, tags.String())
		})
	})
}
