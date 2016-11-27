package render

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

type Engine struct {
	*Options
	moot *sync.Mutex
}

func New(opts *Options) *Engine {
	if opts.TemplateFuncs == nil {
		opts.TemplateFuncs = DefaultHelpers
	}
	helpers := template.FuncMap{}
	for k, v := range DefaultHelpers {
		helpers[k] = v
	}
	for k, v := range opts.TemplateFuncs {
		helpers[k] = v
	}
	opts.TemplateFuncs = helpers

	opts.templates = template.New("").Funcs(helpers)
	if opts.TemplatesPath != "" {
		var err error
		opts.templates, err = parseAndCache(opts.TemplatesPath, helpers)
		if err != nil {
			log.Fatal(err)
		}
	}

	e := &Engine{
		Options: opts,
		moot:    &sync.Mutex{},
	}
	return e
}

var DefaultHelpers = template.FuncMap{
	"yield": func() template.HTML {
		return template.HTML("")
	},
	"partial": func(name string) template.HTML {
		return template.HTML("")
	},
	"debug": func(data interface{}) template.HTML {
		return template.HTML(fmt.Sprintf("%+v", data))
	},
}

func parseAndCache(templatesPath string, helpers template.FuncMap) (*template.Template, error) {
	t := template.New("").Funcs(helpers)
	err := filepath.Walk(templatesPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.WithStack(err)
			}
			t, err = t.New(info.Name()).Parse(string(b))
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})
	return t, errors.WithStack(err)
}
