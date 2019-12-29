package core

import (
	"testing"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	app := meta.New(".")
	app.Name = name.New("buffalo")

	opts := &Options{
		App: app,
	}

	err := opts.Validate()
	r.Error(err)

	opts.App.Name = name.New("coke")
	err = opts.Validate()
	r.NoError(err)

	opts.App.Name = name.New("#$(@#)")
	err = opts.Validate()
	r.Error(err)

	opts.App.Name = name.New("coke")
	err = opts.Validate()
	r.NoError(err)

	opts.App.Name = name.New("test")
	err = opts.Validate()
	r.Error(err)

	opts.App.Name = name.New("testapp")
	err = opts.Validate()
	r.NoError(err)
}
