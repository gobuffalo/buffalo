package meta

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New_Defaults(t *testing.T) {
	r := require.New(t)

	app := New("")
	r.NotEmpty(app.Root)
	r.False(app.WithPop)

	app = New(".")
	r.NotEmpty(app.Root)
}

func Test_New_With_DatabaseYML(t *testing.T) {
	r := require.New(t)

	tmp := t.TempDir()
	dbYML := filepath.Join(tmp, "database.yml")
	r.NoError(os.WriteFile(dbYML, []byte("test"), 0644))

	app := New(tmp)
	r.Equal(tmp, app.Root)
	r.True(app.WithPop)
}

func Test_New_With_TOML(t *testing.T) {
	r := require.New(t)

	tmp := t.TempDir()
	configDir := filepath.Join(tmp, "config")
	r.NoError(os.MkdirAll(configDir, 0755))

	tomlContent := `with_pop = true`
	tomlPath := filepath.Join(configDir, "buffalo-app.toml")
	r.NoError(os.WriteFile(tomlPath, []byte(tomlContent), 0644))

	app := New(tmp)
	r.Equal(tmp, app.Root)
	r.True(app.WithPop)
}

func Test_New_TOML_Priority_Over_DatabaseYML(t *testing.T) {
	r := require.New(t)

	tmp := t.TempDir()

	// Create both files
	configDir := filepath.Join(tmp, "config")
	r.NoError(os.MkdirAll(configDir, 0755))

	tomlContent := `with_pop = false`
	tomlPath := filepath.Join(configDir, "buffalo-app.toml")
	r.NoError(os.WriteFile(tomlPath, []byte(tomlContent), 0644))

	dbYML := filepath.Join(tmp, "database.yml")
	r.NoError(os.WriteFile(dbYML, []byte("test"), 0644))

	// TOML should take priority
	app := New(tmp)
	r.False(app.WithPop)
}

func Test_New_No_Files(t *testing.T) {
	r := require.New(t)

	tmp := t.TempDir()

	app := New(tmp)
	r.Equal(tmp, app.Root)
	r.False(app.WithPop)
}
