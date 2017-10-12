package build

import (
	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/envy"
	pack "github.com/gobuffalo/packr/builder"
	"github.com/pkg/errors"
)

func (b *Builder) buildAssets() error {
	if b.WithWebpack {
		envy.Set("NODE_ENV", "production")
		err := b.exec(webpack.BinPath)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	p := pack.New(b.ctx, b.Root)
	p.Compress = b.Compress

	if b.ExtractAssets {
		p.IgnoredBoxes = append(p.IgnoredBoxes, "../public/assets")
		err := b.buildExtractedAssets()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return p.Run()
}
