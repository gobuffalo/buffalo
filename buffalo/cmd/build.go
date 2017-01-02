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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/velvet"
	"github.com/spf13/cobra"
)

var output string

var cleanup = []string{}

func buildWebpack() error {
	_, err := os.Stat("webpack.config.js")
	if err == nil {
		// build webpack
		cmd := exec.Command("webpack")
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	return nil
}

func buildAPack() error {
	cleanup = append(cleanup, "a")
	err := os.MkdirAll("a", 0766)
	if err != nil {
		return err
	}
	err = buildAInit()
	if err != nil {
		return err
	}
	err = buildDatabase()
	if err != nil {
		return err
	}
	return nil
}

func buildAInit() error {
	path := filepath.Join("a", "a.go")
	cleanup = append(cleanup, path)
	a, err := os.Create(path)
	if err != nil {
		return err
	}
	a.WriteString(`package a

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
}`)
	return nil
}

func buildDatabase() error {
	bb := &bytes.Buffer{}
	path := filepath.Join("a", "database.go")
	dgo, err := os.Create(path)
	cleanup = append(cleanup, path)
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

func buildRice() error {
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)
	_, err := exec.LookPath("rice")
	if err == nil {
		// if rice exists, try and build some cleanup:
		err = filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if filepath.Base(path) == "node_modules" {
					return filepath.SkipDir
				}
				err = os.Chdir(path)
				if err != nil {
					return err
				}
				cmd := exec.Command("rice", "embed-go")
				err = cmd.Run()
				if err == nil {
					bp := filepath.Join(path, "rice-box.go")
					_, err := os.Stat(bp)
					if err == nil {
						fmt.Printf("--> built rice box %s\n", bp)
						cleanup = append(cleanup, bp)
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func buildMain() error {
	root, err := rootPath("")
	if err != nil {
		return err
	}
	ctx := velvet.NewContext()
	ctx.Set("root", root)
	ctx.Set("modelsPack", filepath.Join(packagePath(root), "models"))
	ctx.Set("aPack", filepath.Join(packagePath(root), "a"))
	ctx.Set("name", filepath.Base(root))
	s, err := velvet.Render(buildMainTmpl, ctx)
	if err != nil {
		return err
	}
	path := "buffalo_build_main.go"
	cleanup = append(cleanup, path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	f.WriteString(s)

	return nil
}

func cleanupBuild(original_main []byte) {
	fmt.Println("--> cleaning up build")
	for _, b := range cleanup {
		fmt.Printf("--> cleaning up %s\n", b)
		os.RemoveAll(b)
	}
	maingo, _ := os.Create("main.go")
	maingo.Write(original_main)
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds a Buffalo binary, including bundling of assets (go.rice & webpack)",
	RunE: func(cc *cobra.Command, args []string) error {
		original_main := &bytes.Buffer{}
		maingo, err := os.Open("main.go")
		_, err = original_main.ReadFrom(maingo)
		if err != nil {
			return err
		}
		maingo.Close()
		defer cleanupBuild(original_main.Bytes())

		new_main := strings.Replace(original_main.String(), "func main()", "func original_main()", 1)
		maingo, err = os.Create("main.go")
		if err != nil {
			return err
		}
		_, err = maingo.WriteString(new_main)
		if err != nil {
			return err
		}

		err = buildWebpack()
		if err != nil {
			return err
		}

		err = buildAPack()
		if err != nil {
			return err
		}

		err = buildMain()
		if err != nil {
			return err
		}

		err = buildRice()
		if err != nil {
			return err
		}
		// go build -ldflags "-X main.build =$(git rev-parse --short HEAD)"

		buildArgs := []string{"build", "-v", "-o", output}
		_, err = exec.LookPath("git")
		version := fmt.Sprintf("\"%s\"", time.Now().Format(time.RFC3339))
		if err == nil {
			cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
			out := &bytes.Buffer{}
			cmd.Stdout = out
			err = cmd.Run()
			if err == nil && out.String() != "" {
				version = out.String()
			}
		}
		buildArgs = append(buildArgs, "-ldflags", fmt.Sprintf("-X main.version=%s", version))

		cmd := exec.Command("go", buildArgs...)
		fmt.Printf("--> building %s\n", strings.Join(cmd.Args, " "))
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
	pwd, _ := os.Getwd()
	buildCmd.Flags().StringVarP(&output, "output", "o", filepath.Join("bin", filepath.Base(pwd)), "set the name of the binary")
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
var migrationBox *rice.Box

func main() {
	fmt.Printf("{{name}} version %s\n\n", version)
	args := os.Args
	if len(args) == 1 {
		original_main()
	}
	c := args[1]
	switch c {
	case "migrate":
		migrate()
	case "start", "run", "serve":
		original_main()
	default:
		log.Fatalf("Could not find a command named: %s", c)
	}
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
