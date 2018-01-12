package build

import (
	"context"
	"os"

	"github.com/gobuffalo/packr/builder"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Builder builds a Buffalo binary
type Builder struct {
	Options
	ctx       context.Context
	steps     []func() error
	cleanups  []func() error
	originals map[string][]byte
}

// New Builder
func New(ctx context.Context, opts Options) *Builder {
	b := &Builder{
		ctx:       ctx,
		Options:   opts,
		cleanups:  []func() error{},
		originals: map[string][]byte{},
	}

	b.steps = []func() error{
		b.prepTarget,
		b.transformMain,
		b.createBuildMain,
		b.prepAPackage,
		b.buildAInit,
		b.buildADatabase,
		b.buildAssets,
		b.buildBin,
	}

	return b
}

// Run builds a Buffalo binary
func (b *Builder) Run() error {
	defer b.Cleanup()
	logrus.Debug(b.Options)

	for _, s := range b.steps {
		err := s()
		if err != nil {
			return errors.WithStack(err)
		}
		os.Chdir(b.Root)
	}
	return nil
}

// Cleanup after run. This is automatically run at the end of `Run`.
// It is provided publicly in-case a manual clean up is necessary.
func (b *Builder) Cleanup() error {
	builder.Clean(b.Root)
	me := multiError{}
	for _, c := range b.cleanups {
		err := c()
		if err != nil {
			me = append(me, errors.Errorf("error while cleaning up: %+v", err))
		}
	}
	if len(me) > 0 {
		logrus.Error(me)
		return me
	}
	return nil
}
