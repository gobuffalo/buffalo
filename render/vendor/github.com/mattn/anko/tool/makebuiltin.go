package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func pkgName(f string) string {
	file, err := parser.ParseFile(token.NewFileSet(), f, nil, parser.PackageClauseOnly)
	if err != nil || file == nil {
		return ""
	}
	return file.Name.Name
}

func isGoFile(dir os.FileInfo) bool {
	return !dir.IsDir() &&
		!strings.HasPrefix(dir.Name(), ".") && // ignore .files
		filepath.Ext(dir.Name()) == ".go"
}

func isPkgFile(dir os.FileInfo) bool {
	return isGoFile(dir) && !strings.HasSuffix(dir.Name(), "_test.go") // ignore test files
}

func parseDir(p string) (map[string]*ast.Package, error) {
	_, pn := filepath.Split(p)

	isGoDir := func(d os.FileInfo) bool {
		if isPkgFile(d) {
			name := pkgName(p + "/" + d.Name())
			return name == pn
		}
		return false
	}

	pkgs, err := parser.ParseDir(token.NewFileSet(), p, isGoDir, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

func main() {
	pkg := "flag"
	if len(os.Args) == 2 {
		pkg = os.Args[1]
	}
	b, err := exec.Command("go", "env", "GOROOT").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	paths := []string{filepath.Join(strings.TrimSpace(string(b)), "src")}
	b, err = exec.Command("go", "env", "GOPATH").CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range strings.Split(strings.TrimSpace(string(b)), string(filepath.ListSeparator)) {
		paths = append(paths, filepath.Join(p, "src"))
	}
	for _, p := range paths {
		pp := filepath.Join(p, pkg)
		pkgs, err := parseDir(pp)
		if err != nil {
			continue
		}
		names := map[string]bool{}
		for _, pp := range pkgs {
			for _, f := range pp.Files {
				for _, d := range f.Decls {
					switch decl := d.(type) {
					case *ast.GenDecl:
						for _, spec := range decl.Specs {
							if vspec, ok := spec.(*ast.ValueSpec); ok {
								for _, n := range vspec.Names {
									c := n.Name[0]
									if c < 'A' || c > 'Z' {
										continue
									}
									names[n.Name] = true
								}
							}
						}
					case *ast.FuncDecl:
						if decl.Recv != nil {
							continue
						}
						c := decl.Name.Name[0]
						if c < 'A' || c > 'Z' {
							continue
						}
						names[decl.Name.Name] = true
					}
				}
			}
		}
		keys := []string{}
		for k, _ := range names {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		_, pn := filepath.Split(pkg)
		fmt.Printf(`// Package %s implements %s interface for anko script.
package %s

import (
	"github.com/mattn/anko/vm"
	pkg "%s"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewModule("%s")
`, pn, pkg, pn, pkg, pn)
		for _, k := range keys {
			fmt.Printf("\t"+`m.Define("%s", pkg.%s)`+"\n", k, k)
		}
		fmt.Println("\treturn m")
		fmt.Println("}")
	}
}
