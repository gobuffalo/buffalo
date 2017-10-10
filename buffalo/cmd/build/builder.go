package build

import (
	"context"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr/builder"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	Options
	ctx       context.Context
	steps     []func() error
	cleanups  []func() error
	originals map[string][]byte
}

func New(ctx context.Context, opts Options) *Builder {
	if _, err := os.Stat(filepath.Join(opts.Root, "database.yml")); err == nil {
		opts.HasDB = true
	}
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
