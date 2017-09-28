package cmd

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/envy"
	pack "github.com/gobuffalo/packr/builder"
	"github.com/gobuffalo/plush"
	"github.com/markbates/sigtx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var outputBinName string
var extractAssets bool
var hasDB bool
var ldflags string
var buildTags string
var static bool
var compress bool

type builder struct {
	cleanup      []string
	originalMain []byte
	originalApp  []byte
	workDir      string
	buildTags    []string
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
		return errors.WithStack(err)
	}
	err = b.buildAInit()
	if err != nil {
		return errors.WithStack(err)
	}
	err = b.buildDatabase()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (b *builder) buildAInit() error {
	a, err := os.Create(b.clean("a", "a.go"))
	if err != nil {
		return errors.WithStack(err)
	}
	a.WriteString(aGo)
	return nil
}

func (b *builder) buildDatabase() error {
	bb := &bytes.Buffer{}
	dgo, err := os.Create(b.clean("a", "database.go"))
	if err != nil {
		return errors.WithStack(err)
	}
	if hasDB {
		// copy the database.yml file to the migrations folder so it's available through packr
		os.MkdirAll("./migrations", 0755)
		d, err := os.Open("database.yml")
		if err != nil {
			return errors.WithStack(err)
		}
		_, err = io.Copy(bb, d)
		if err != nil {
			return errors.WithStack(err)
		}
		if !bytes.Contains(bb.Bytes(), []byte("sqlite")) {
			b.buildTags = append(b.buildTags, "nosqlite")
		} else if !static {
			fmt.Println("you are building a SQLite application, please consider using the `--static` flag to compile a static binary")
		}

	} else {
		// add the nosqlite build tag if there is no database being used
		b.buildTags = append(b.buildTags, "nosqlite")
	}
	dgo.WriteString("package a\n")
	dgo.WriteString(fmt.Sprintf("var DB_CONFIG = `%s`", bb.String()))
	return nil
}

func (b *builder) buildPackrEmbedded() error {
	defer os.Chdir(b.workDir)
	p := pack.New(context.Background(), b.workDir)
	p.Compress = compress
	if extractAssets {
		p.IgnoredBoxes = append(p.IgnoredBoxes, "../public/assets")
	}
	return p.Run()
}

func (b *builder) disableAssetsHandling() error {
	defer os.Chdir(b.workDir)
	fmt.Printf("--> disable self assets handling\n")

	newApp := strings.Replace(string(b.originalApp), "app.ServeFiles(\"/assets\", assetsPath())", "//app.ServeFiles(\"/assets\", assetsPath())", 1)

	appgo, err := os.Create("actions/app.go")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = appgo.WriteString(newApp)
	if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return errors.WithStack(err)
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return errors.WithStack(err)
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
			return errors.WithStack(err)
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return errors.WithStack(err)
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return errors.WithStack(err)
	})

	return errors.WithStack(err)
}

func (b *builder) buildMain() error {
	newMain := strings.Replace(string(b.originalMain), "func main()", "func originalMain()", 1)
	maingo, err := os.Create("main.go")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = maingo.WriteString(newMain)
	if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}
	f, err := os.Create(b.clean("buffalo_build_main.go"))
	if err != nil {
		return errors.WithStack(err)
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

	pack.Clean(b.workDir)

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

}

func (b *builder) run() error {
	_, err := os.Stat("database.yml")
	if err == nil {
		hasDB = true
	}

	err = b.buildMain()
	if err != nil {
		return errors.WithStack(err)
	}

	err = b.buildWebpack()
	if err != nil {
		return errors.WithStack(err)
	}

	err = b.buildAPack()
	if err != nil {
		return errors.WithStack(err)
	}

	err = b.buildMain()
	if err != nil {
		return errors.WithStack(err)
	}

	if extractAssets {
		err = b.buildAssetsArchive()
		if err != nil {
			return errors.WithStack(err)
		}
		err = b.disableAssetsHandling()
		if err != nil {
			return errors.WithStack(err)
		}
	}

	err = b.buildPackrEmbedded()
	if err != nil {
		return errors.WithStack(err)
	}
	return b.buildBin()
}

func (b *builder) buildBin() error {
	buildArgs := []string{"build", "-v"}
	if len(b.buildTags) > 0 {
		buildArgs = append(buildArgs, "-tags", strings.Join(b.buildTags, " "))
	}
	buildArgs = append(buildArgs, "-o", outputBinName)
	_, err := exec.LookPath("git")
	buildTime := fmt.Sprintf("\"%s\"", time.Now().Format(time.RFC3339))
	version := buildTime
	if err == nil {
		cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
		out := &bytes.Buffer{}
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err == nil && out.String() != "" {
			version = strings.TrimSpace(out.String())
		}
	}

	// RUN GOOS=linux buffalo build --ldflags '-linkmode external -extldflags "-static"' -o /bin/app
	flags := []string{
		fmt.Sprintf("-X main.version=%s", version),
		fmt.Sprintf("-X main.buildTime=%s", buildTime),
	}
	if static {
		flags = append(flags, "-linkmode external", "-extldflags \"-static\"")
	}

	// Add any additional ldflags passed in to the build args
	if len(ldflags) > 0 {
		if foundVersion, _ := regexp.MatchString("-X\\s+main.version=", ldflags); foundVersion {
			return errors.New("the ldflag option '-X main.version=' is reserved for Buffalo use")
		}
		if foundBuildTime, _ := regexp.MatchString("-X\\s+main.buildTime=", ldflags); foundBuildTime {
			return errors.New("the ldflag option '-X main.buildTime=' is reserved for Buffalo use")
		}
		flags = append(flags, ldflags)
	}

	buildArgs = append(buildArgs, "-ldflags", strings.Join(flags, " "))

	return b.exec(envy.Get("GO_BIN", "go"), buildArgs...)
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b", "bill"},
	Short:   "Builds a Buffalo binary, including bundling of assets (packr & webpack)",
	RunE: func(cc *cobra.Command, args []string) error {
		originalMain := &bytes.Buffer{}
		maingo, err := os.Open("main.go")
		_, err = originalMain.ReadFrom(maingo)
		if err != nil {
			return errors.WithStack(err)
		}
		maingo.Close()

		originalApp := &bytes.Buffer{}
		appgo, err := os.Open("actions/app.go")
		_, err = originalApp.ReadFrom(appgo)
		if err != nil {
			return errors.WithStack(err)
		}
		appgo.Close()

		pwd, _ := os.Getwd()
		b := builder{
			cleanup:      []string{},
			originalMain: originalMain.Bytes(),
			originalApp:  originalApp.Bytes(),
			workDir:      pwd,
			buildTags:    []string{},
		}
		if buildTags != "" {
			b.buildTags = append(b.buildTags, buildTags)
		}

		ctx, cancel := sigtx.WithCancel(context.Background(), os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)
		go func() {
			<-ctx.Done()
			if ctx.Err() == context.Canceled {
				fmt.Println("~~~BUILD CANCELLED ~~~")
				b.cleanupBuild()
			}
		}()
		defer cancel()

		defer b.cleanupBuild()
		b.cleanupTarget()
		return b.run()
	},
}

func init() {
	decorate("build", buildCmd)

	RootCmd.AddCommand(buildCmd)
	pwd, _ := os.Getwd()
	output := filepath.Join("bin", filepath.Base(pwd))

	if runtime.GOOS == "windows" {
		output += ".exe"
	}

	buildCmd.Flags().StringVarP(&outputBinName, "output", "o", output, "set the name of the binary")
	buildCmd.Flags().StringVarP(&buildTags, "tags", "t", "", "compile with specific build tags")
	buildCmd.Flags().BoolVarP(&extractAssets, "extract-assets", "e", false, "extract the assets and put them in a distinct archive")
	buildCmd.Flags().BoolVarP(&static, "static", "s", false, "build a static binary using  --ldflags '-linkmode external -extldflags \"-static\"' (USE FOR CGO)")
	buildCmd.Flags().BoolVarP(&compress, "compress", "c", true, "compress static files in the binary")
	buildCmd.Flags().StringVar(&ldflags, "ldflags", "", "set any ldflags to be passed to the go build")
}

var buildMainTmpl = `package main

import (
	"fmt"
	"log"
	"os"

	"github.com/markbates/grift/grift"
	_ "<%= aPack %>"
	<%= if (modelsPack) { %>
	"github.com/gobuffalo/packr"
	"github.com/markbates/pop"
	"<%= modelsPack %>"
	<% } %>
	<%= if (griftsPack) { %>
	_ "<%= griftsPack %>"
	<% } %>
)

var version = "unknown"
var buildTime = "unknown"

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
		c := grift.NewContext(args[2])
		if len(args) > 2 {
			c.Args = args[3:]
		}
		err := grift.Run(args[2], c)
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
	box, err := pop.NewMigrationBox(packr.NewBox("./migrations"), models.DB)
	if err != nil {
		log.Fatalf("Failed to unpack migrations: %s", err)
	}
	err = box.Up()
	if err != nil {
		log.Fatalf("Failed to run migrations: %s", err)
	}
}
<% } %>
`

var aGo = `package a

import (
	"log"
	"strings"
	"github.com/markbates/pop"
)

func init() {
	dropDatabaseYml()
}

func dropDatabaseYml() {
	if DB_CONFIG != "" {
		r := strings.NewReader(DB_CONFIG)
		err := pop.LoadFrom(r)
		if err != nil {
			log.Fatal(err)
		}
	}
}`
