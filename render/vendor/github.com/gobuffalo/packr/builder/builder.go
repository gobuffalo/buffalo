package builder

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/pkg/errors"
)

var boxPattern = regexp.MustCompile(`packr.NewBox\(["` + "`" + `]([^\(\)]+)["` + "`" + `]\)`)

var packagePattern = regexp.MustCompile(`package\s+(\w+)`)

//var packagePattern = regexp.MustCompile(`package\s+([^\r]+)`)
var invalidFilePattern = regexp.MustCompile(`(_test|-packr).go$`)

// Builder scans folders/files looking for `packr.NewBox` and then compiling
// the required static files into `<package-name>-packr.go` files so they can
// be built into Go binaries.
type Builder struct {
	context.Context
	RootPath string
	pkgs     map[string]pkg
}

// Run the builder.
func (b *Builder) Run() error {
	err := filepath.Walk(b.RootPath, func(path string, info os.FileInfo, err error) error {
		base := filepath.Base(path)
		if base == ".git" || base == "vendor" || base == "node_modules" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			return b.process(path)
		}
		return nil
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return b.dump()
}

func (b *Builder) dump() error {
	for _, p := range b.pkgs {
		name := filepath.Join(p.Dir, p.Name+"-packr.go")
		fmt.Printf("--> packing %s\n", name)
		f, err := os.Create(name)
		defer f.Close()
		if err != nil {
			return errors.WithStack(err)
		}
		t, err := template.New("").Parse(tmpl)

		if err != nil {
			return errors.WithStack(err)
		}
		err = t.Execute(f, p)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func (b *Builder) process(path string) error {
	ext := filepath.Ext(path)
	if ext != ".go" || invalidFilePattern.MatchString(path) {
		return nil
	}

	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.WithStack(err)
	}
	fb := string(bb)

	matches := boxPattern.FindAllStringSubmatch(fb, -1)
	if len(matches) == 0 {
		return nil
	}

	pk := pkg{
		Dir:   filepath.Dir(path),
		Boxes: []box{},
	}
	pname := packagePattern.FindStringSubmatch(fb)
	pk.Name = pname[1]

	for _, m := range matches {
		bx := &box{
			Name:  m[1],
			Files: []file{},
		}
		err = bx.Walk(filepath.Join(pk.Dir, bx.Name))
		if err != nil {
			return errors.WithStack(err)
		}
		if len(bx.Files) > 0 {
			pk.Boxes = append(pk.Boxes, *bx)
		}
	}

	if len(pk.Boxes) > 0 {
		b.addPkg(pk)
	}
	return nil
}

func (b *Builder) addPkg(p pkg) {
	if _, ok := b.pkgs[p.Name]; !ok {
		b.pkgs[p.Name] = p
		return
	}
	pp := b.pkgs[p.Name]
	pp.Boxes = append(pp.Boxes, p.Boxes...)
	b.pkgs[p.Name] = pp
}

// New Builder with a given context and path
func New(ctx context.Context, path string) *Builder {
	return &Builder{
		Context:  ctx,
		RootPath: path,
		pkgs:     map[string]pkg{},
	}
}
