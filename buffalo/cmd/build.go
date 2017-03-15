package cmd

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	pack "github.com/gobuffalo/packr/builder"
	"github.com/gobuffalo/plush"
	"github.com/spf13/cobra"
)

var outputBinName string
var zipBin bool
var extractAssets bool
var hasDB bool

type builder struct {
	cleanup      []string
	originalMain []byte
	originalApp  []byte
	workDir      string
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
		return b.exec(webpack.BinPath)
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
	if hasDB {
		// copy the database.yml file to the migrations folder so it's available through rice
		os.MkdirAll("./migrations", 0755)
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

// func (b *builder) buildRiceZip() error {
// 	defer os.Chdir(b.workDir)
// 	_, err := exec.LookPath("rice")
// 	if err == nil {
// 		paths := map[string]bool{}
// 		// if rice exists, try and build some cleanup:
// 		err = filepath.Walk(b.workDir, func(path string, info os.FileInfo, err error) error {
// 			if info.IsDir() {
// 				base := filepath.Base(path)
// 				if base == "node_modules" || base == ".git" || base == "bin" || base == "vendor" {
// 					return filepath.SkipDir
// 				}
// 			} else {
// 				err = os.Chdir(filepath.Dir(path))
// 				if err != nil {
// 					return err
// 				}
//
// 				s, err := ioutil.ReadFile(path)
// 				if err != nil {
// 					return err
// 				}
// 				rx := regexp.MustCompile("(rice.FindBox|rice.MustFindBox)")
// 				if rx.Match(s) && filepath.Ext(info.Name()) == ".go" {
// 					gopath := strings.Replace(filepath.Join(os.Getenv("GOPATH"), "src"), "\\", "/", -1)
// 					pkg := strings.Replace(filepath.Dir(strings.Replace(path, gopath+"/", "", -1)), "\\", "/", -1)
// 					paths[pkg] = true
// 				}
// 			}
// 			return nil
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		if len(paths) != 0 {
// 			args := []string{"append", "--exec", filepath.Join(b.workDir, outputBinName)}
// 			for k := range paths {
// 				args = append(args, "-i", k)
// 			}
// 			return b.exec("rice", args...)
// 		}
// 		// rice append --exec example
// 	}
// 	return nil
// }
func (b *builder) buildPackrEmbedded() error {
	defer os.Chdir(b.workDir)
	defer pack.Clean(b.workDir)
	p := pack.New(context.Background(), b.workDir)
	return p.Run()
}

func (b *builder) disableAssetsHandling() error {
	defer os.Chdir(b.workDir)
	fmt.Printf("--> disable self assets handling\n")

	newApp := strings.Replace(string(b.originalApp), "app.ServeFiles(\"/assets\", assetsPath())", "//app.ServeFiles(\"/assets\", assetsPath())", 1)

	appgo, err := os.Create("actions/app.go")
	if err != nil {
		return err
	}
	_, err = appgo.WriteString(newApp)
	if err != nil {
		return err
	}

	return nil
}

func (b *builder) buildAssetsArchive() error {
	defer os.Chdir(b.workDir)
	fmt.Printf("--> build assets archive\n")

	outputDir := filepath.Dir(outputBinName)
	assetsName := filepath.Base(outputBinName)
	target := outputDir + "/" + assetsName + "-assets.zip"
	source := filepath.Join(b.workDir, "public", "assets")

	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return err
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func (b *builder) buildMain() error {
	newMain := strings.Replace(string(b.originalMain), "func main()", "func originalMain()", 1)
	maingo, err := os.Create("main.go")
	if err != nil {
		return err
	}
	_, err = maingo.WriteString(newMain)
	if err != nil {
		return err
	}

	ctx := plush.NewContext()
	ctx.Set("root", rootPath)
	ctx.Set("hasDB", hasDB)
	if hasDB {
		ctx.Set("modelsPack", packagePath(rootPath)+"/models")
	}
	_, err = os.Stat(filepath.Join(rootPath, "grifts"))
	if err == nil {
		ctx.Set("griftsPack", packagePath(rootPath)+"/grifts")
	}
	ctx.Set("aPack", packagePath(rootPath)+"/a")
	ctx.Set("name", filepath.Base(rootPath))
	s, err := plush.Render(buildMainTmpl, ctx)
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
		fmt.Printf("----> cleaning up %s\n", b)
		os.RemoveAll(b)
	}
	maingo, _ := os.Create("main.go")
	maingo.Write(b.originalMain)

	appgo, _ := os.Create("actions/app.go")
	appgo.Write(b.originalApp)
}

func (b *builder) cleanupTarget() {
	fmt.Println("--> cleaning up target dir")

	// Create output directory if not exists
	outputDir := filepath.Dir(outputBinName)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0776)
		fmt.Printf("----> creating target dir %s\n", outputDir)
	}

	files, _ := ioutil.ReadDir(outputDir)
	for _, f := range files {
		fmt.Printf("----> cleaning up %s\n", f.Name())
		os.RemoveAll(outputDir + f.Name())
	}
}

func (b *builder) run() error {
	_, err := os.Stat("database.yml")
	if err == nil {
		hasDB = true
	}

	err = b.buildMain()
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

	if extractAssets {
		err = b.buildAssetsArchive()
		if err != nil {
			return err
		}
		err = b.disableAssetsHandling()
		if err != nil {
			return err
		}
		return b.buildBin()
	}

	// if zipBin {
	// 	err = b.buildBin()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return b.buildRiceZip()
	// }

	err = b.buildPackrEmbedded()
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
		originalMain := &bytes.Buffer{}
		maingo, err := os.Open("main.go")
		_, err = originalMain.ReadFrom(maingo)
		if err != nil {
			return err
		}
		maingo.Close()

		originalApp := &bytes.Buffer{}
		appgo, err := os.Open("actions/app.go")
		_, err = originalApp.ReadFrom(appgo)
		if err != nil {
			return err
		}
		appgo.Close()

		pwd, _ := os.Getwd()
		b := builder{
			cleanup:      []string{},
			originalMain: originalMain.Bytes(),
			originalApp:  originalApp.Bytes(),
			workDir:      pwd,
		}
		defer b.cleanupBuild()

		b.cleanupTarget()
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
	buildCmd.Flags().BoolVarP(&extractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
}

var buildMainTmpl = `package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markbates/grift/grift"
	rice "github.com/GeertJohan/go.rice"
	_ "<%= aPack %>"
	<%= if (modelsPack) { %>
	"io/ioutil"
	"path/filepath"
	"<%= modelsPack %>"
	<% } %>
	<%= if (griftsPack) { %>
	_ "<%= griftsPack %>"
	<% } %>
)

var version = "unknown"
var buildTime = "unknown"
var migrationBox *rice.Box

func main() {
	args := os.Args
	if len(args) == 1 {
		originalMain()
	}
	c := args[1]
	switch c {
	<%= if (modelsPack) { %>
	case "migrate":
		migrate()
	<% } %>
	case "start", "run", "serve":
		printVersion()
		originalMain()
	case "version":
		printVersion()
	case "task", "t", "tasks":
		err := grift.Run(args[2], grift.NewContext(args[2]))
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Could not find a command named: %s", c)
	}
}

func printVersion() {
	fmt.Printf("<%= name %> version %s (%s)\n\n", version, buildTime)
}

<%= if (modelsPack) { %>
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
	dir, err := ioutil.TempDir("", "<%= name %>-migrations")
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
}
<% } %>
`

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
