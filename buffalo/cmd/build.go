// Copyright © 2016 Mark Bates <mark@markbates.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gobuffalo/velvet"
	"github.com/spf13/cobra"
)

var outputBinName string
var zipBin bool

type builder struct {
	cleanup       []string
	original_main []byte
	workDir       string
}

func (b *builder) clean(name ...string) string {
	path := filepath.Join(name...)
	b.cleanup = append(b.cleanup, path)
	return path
}

func (b *builder) exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	fmt.Printf("--> running %s\n", strings.Join(cmd.Args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (b *builder) execQuiet(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func (b *builder) buildWebpack() error {
	_, err := os.Stat("webpack.config.js")
	if err == nil {
		// build webpack
		return b.exec("webpack")
	}
	return nil
}

func (b *builder) buildAPack() error {
	err := os.MkdirAll(b.clean("a"), 0766)
	if err != nil {
		return err
	}
	err = b.buildAInit()
	if err != nil {
		return err
	}
	err = b.buildDatabase()
	if err != nil {
		return err
	}
	return nil
}

func (b *builder) buildAInit() error {
	a, err := os.Create(b.clean("a", "a.go"))
	if err != nil {
		return err
	}
	a.WriteString(aGo)
	return nil
}

func (b *builder) buildDatabase() error {
	bb := &bytes.Buffer{}
	dgo, err := os.Create(b.clean("a", "database.go"))
	if err != nil {
		return err
	}
	_, err = os.Stat("database.yml")
	if err == nil {
		// copy the database.yml file to the migrations folder so it's available through rice
		d, err := os.Open("database.yml")
		if err != nil {
			return err
		}
		_, err = io.Copy(bb, d)
		if err != nil {
			return err
		}
	}
	dgo.WriteString("package a\n")
	dgo.WriteString(fmt.Sprintf("var DB_CONFIG = `%s`", bb.String()))
	return nil
}

func (b *builder) buildRiceZip() error {
	defer os.Chdir(b.workDir)
	_, err := exec.LookPath("rice")
	if err == nil {
		paths := map[string]bool{}
		// if rice exists, try and build some cleanup:
		err = filepath.Walk(b.workDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				base := filepath.Base(path)
				if base == "node_modules" || base == ".git" || base == "bin" {
					return filepath.SkipDir
				}
			} else {
				err = os.Chdir(filepath.Dir(path))
				if err != nil {
					return err
				}

				s, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				rx := regexp.MustCompile("(rice.FindBox|rice.MustFindBox)")
				if rx.Match(s) {
					gopath := strings.Replace(filepath.Join(os.Getenv("GOPATH"), "src"), "\\", "/", -1)
					pkg := strings.Replace(filepath.Dir(strings.Replace(path, gopath+"/", "", -1)), "\\", "/", -1)
					paths[pkg] = true
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		if len(paths) != 0 {
			args := []string{"append", "--exec", filepath.Join(b.workDir, outputBinName)}
			for k := range paths {
				args = append(args, "-i", k)
			}
			return b.exec("rice", args...)
		}
		// rice append --exec example
	}
	return nil
}
func (b *builder) buildRiceEmbedded() error {
	defer os.Chdir(b.workDir)
	_, err := exec.LookPath("rice")
	if err == nil {
		// if rice exists, try and build some cleanup:
		err = filepath.Walk(b.workDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				base := filepath.Base(path)
				if base == "node_modules" || base == ".git" {
					return filepath.SkipDir
				}
				err = os.Chdir(path)
				if err != nil {
					return err
				}
				err = b.execQuiet("rice", "embed-go")
				if err == nil {
					bp := filepath.Join(path, "rice-box.go")
					_, err := os.Stat(bp)
					if err == nil {
						fmt.Printf("--> built rice box %s\n", bp)
						b.clean(bp)
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		// rice append --exec example
	}
	return nil
}

func (b *builder) buildMain() error {
	new_main := strings.Replace(string(b.original_main), "func main()", "func original_main()", 1)
	maingo, err := os.Create("main.go")
	if err != nil {
		return err
	}
	_, err = maingo.WriteString(new_main)
	if err != nil {
		return err
	}

	root, err := rootPath("")
	if err != nil {
		return err
	}
	ctx := velvet.NewContext()
	ctx.Set("root", root)
	ctx.Set("modelsPack", packagePath(root)+"/models")
	ctx.Set("aPack", packagePath(root)+"/a")
	ctx.Set("name", filepath.Base(root))
	s, err := velvet.Render(buildMainTmpl, ctx)
	if err != nil {
		return err
	}
	f, err := os.Create(b.clean("buffalo_build_main.go"))
	if err != nil {
		return err
	}
	f.WriteString(s)

	return nil
}

func (b *builder) cleanupBuild() {
	fmt.Println("--> cleaning up build")
	for _, b := range b.cleanup {
		fmt.Printf("--> cleaning up %s\n", b)
		os.RemoveAll(b)
	}
	maingo, _ := os.Create("main.go")
	maingo.Write(b.original_main)
}

func (b *builder) run() error {
	err := b.buildMain()
	if err != nil {
		return err
	}

	err = b.buildWebpack()
	if err != nil {
		return err
	}

	err = b.buildAPack()
	if err != nil {
		return err
	}

	err = b.buildMain()
	if err != nil {
		return err
	}

	if zipBin {
		err = b.buildBin()
		if err != nil {
			return err
		}
		return b.buildRiceZip()
	}

	err = b.buildRiceEmbedded()
	if err != nil {
		return err
	}
	return b.buildBin()
}

func (b *builder) buildBin() error {
	buildArgs := []string{"build", "-v", "-o", outputBinName}
	_, err := exec.LookPath("git")
	buildTime := fmt.Sprintf("\"%s\"", time.Now().Format(time.RFC3339))
	version := buildTime
	if err == nil {
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		out := &bytes.Buffer{}
		cmd.Stdout = out
		err = cmd.Run()
		if err == nil && out.String() != "" {
			version = strings.TrimSpace(out.String())
		}
	}
	buildArgs = append(buildArgs, "-ldflags", fmt.Sprintf("-X main.version=%s -X main.buildTime=%s", version, buildTime))

	return b.exec("go", buildArgs...)
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "bill"},
	Short:   "Builds a Buffalo binary, including bundling of assets (go.rice & webpack)",
	RunE: func(cc *cobra.Command, args []string) error {
		original_main := &bytes.Buffer{}
		maingo, err := os.Open("main.go")
		_, err = original_main.ReadFrom(maingo)
		if err != nil {
			return err
		}
		maingo.Close()
		pwd, _ := os.Getwd()
		b := builder{
			cleanup:       []string{},
			original_main: original_main.Bytes(),
			workDir:       pwd,
		}
		defer b.cleanupBuild()

		return b.run()
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	pwd, _ := os.Getwd()
	output := filepath.Join("bin", filepath.Base(pwd))

	if runtime.GOOS == "windows" {
		output += ".exe"
	}

	buildCmd.Flags().StringVarP(&outputBinName, "output", "o", output, "set the name of the binary")
	buildCmd.Flags().BoolVarP(&zipBin, "zip", "z", false, "zips the assets to the binary, this requires zip installed")
}

var buildMainTmpl = `package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	_ "{{aPack}}"
	"{{modelsPack}}"
)

var version = "unknown"
var buildTime = "unknown"
var migrationBox *rice.Box

func main() {
	args := os.Args
	if len(args) == 1 {
		original_main()
	}
	c := args[1]
	switch c {
	case "migrate":
		migrate()
	case "start", "run", "serve":
		printVersion()
		original_main()
	case "version":
		printVersion()
	default:
		log.Fatalf("Could not find a command named: %s", c)
	}
}

func printVersion() {
	fmt.Printf("{{name}} version %s (%s)\n\n", version, buildTime)
}

func migrate() {
	var err error
	migrationBox, err = rice.FindBox("./migrations")
	if err != nil {
		fmt.Println("--> No migrations found.")
		return
	}
	fmt.Println("--> Running migrations")
	path, err := unpackMigrations()
	if err != nil {
		log.Fatalf("Failed to unpack migrations: %s", err)
	}
	defer os.RemoveAll(path)

	models.DB.MigrateUp(path)
}

func unpackMigrations() (string, error) {
	dir, err := ioutil.TempDir("", "{{name}}-migrations")
	if err != nil {
		log.Fatalf("Unable to create temp directory: %s", err)
	}

	migrationBox.Walk("", func(path string, fi os.FileInfo, e error) error {
		if !fi.IsDir() {
			content := migrationBox.MustBytes(path)
			file := filepath.Join(dir, path)
			if err := ioutil.WriteFile(file, content, 0666); err != nil {
				log.Fatalf("Failed to write migration to disk: %s", err)
			}
		}
		return e
	})

	return dir, nil
}`

var aGo = `package a

import (
	"log"
	"os"
)

func init() {
	dropDatabaseYml()
}

func dropDatabaseYml() {
	if DB_CONFIG != "" {

		_, err := os.Stat("database.yml")
		if err == nil {
			// yaml already exists, don't do anything
			return
		}
		f, err := os.Create("database.yml")
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.WriteString(DB_CONFIG)
		if err != nil {
			log.Fatal(err)
		}
	}
}`
