package render

import (
	"html/template"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/tags"
	"github.com/stretchr/testify/require"
)

func Test_javascriptTag(t *testing.T) {
	r := require.New(t)
	re := New(Options{
		AssetsBox: packr.NewBox(""),
	})
	tr := re.Template("").(templateRenderer)
	h := tr.addAssetsHelpers(Helpers{})
	f := h["javascriptTag"].(func(string, tags.Options) (template.HTML, error))
	s, err := f("application.js", nil)
	r.NoError(err)
	r.Equal(template.HTML(`<script src="/assets/application.js" type="text/javascript"></script>`), s)
}

func Test_javascriptTag_Options(t *testing.T) {
	r := require.New(t)
	re := New(Options{
		AssetsBox: packr.NewBox(""),
	})
	tr := re.Template("").(templateRenderer)
	h := tr.addAssetsHelpers(Helpers{})
	f := h["javascriptTag"].(func(string, tags.Options) (template.HTML, error))
	s, err := f("application.js", tags.Options{"class": "foo"})
	r.NoError(err)
	r.Equal(template.HTML(`<script class="foo" src="/assets/application.js" type="text/javascript"></script>`), s)

}
func Test_stylesheetTag(t *testing.T) {
	r := require.New(t)
	re := New(Options{
		AssetsBox: packr.NewBox(""),
	})
	tr := re.Template("").(templateRenderer)
	h := tr.addAssetsHelpers(Helpers{})
	f := h["stylesheetTag"].(func(string, tags.Options) (template.HTML, error))
	s, err := f("application.css", nil)
	r.NoError(err)
	r.Equal(template.HTML(`<link href="/assets/application.css" media="screen" rel="stylesheet" />`), s)
}

func Test_stylesheetTag_Options(t *testing.T) {
	r := require.New(t)
	re := New(Options{
		AssetsBox: packr.NewBox(""),
	})
	tr := re.Template("").(templateRenderer)
	h := tr.addAssetsHelpers(Helpers{})
	f := h["stylesheetTag"].(func(string, tags.Options) (template.HTML, error))
	s, err := f("application.css", tags.Options{"class": "foo"})
	r.NoError(err)
	r.Equal(template.HTML(`<link class="foo" href="/assets/application.css" media="screen" rel="stylesheet" />`), s)
}
func Test_imgTag(t *testing.T) {
	r := require.New(t)
	re := New(Options{
		AssetsBox: packr.NewBox(""),
	})
	tr := re.Template("").(templateRenderer)
	h := tr.addAssetsHelpers(Helpers{})
	f := h["imgTag"].(func(string, tags.Options) (template.HTML, error))
	s, err := f("foo.png", nil)
	r.NoError(err)
	r.Equal(template.HTML(`<img src="/assets/foo.png" />`), s)
}

func Test_imgTag_Options(t *testing.T) {
	r := require.New(t)
	re := New(Options{
		AssetsBox: packr.NewBox(""),
	})
	tr := re.Template("").(templateRenderer)
	h := tr.addAssetsHelpers(Helpers{})
	f := h["imgTag"].(func(string, tags.Options) (template.HTML, error))
	s, err := f("foo.png", tags.Options{"class": "foo"})
	r.NoError(err)
	r.Equal(template.HTML(`<img class="foo" src="/assets/foo.png" />`), s)
}
