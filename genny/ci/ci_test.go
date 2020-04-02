package ci

import (
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Provider: "travis",
		DBType:   "postgres",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal(".travis.yml", f.Name())
	travisYml := struct {
		Language     string
		Go           []string
		Env          []string
		Services     []string
		BeforeScript []string `yaml:"before_script"`
		GoImportPath string   `yaml:"go_import_path"`
		Install      []string
		Script       string
	}{}
	r.NoError(yaml.NewDecoder(f).Decode(&travisYml), ".travis.yml is a valid YAML file")
}

func Test_New_Gitlab(t *testing.T) {
	r := require.New(t)

	app := meta.New(".")
	app.WithPop = true

	g, err := New(&Options{
		App:      app,
		Provider: "gitlab",
		DBType:   "postgres",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal(".gitlab-ci.yml", f.Name())
	r.Contains(f.String(), "postgres:5432")
}

func Test_New_Gitlab_No_pop(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Provider: "gitlab",
		DBType:   "postgres",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal(".gitlab-ci.yml", f.Name())
	r.NotContains(f.String(), "postgres:5432")
}
