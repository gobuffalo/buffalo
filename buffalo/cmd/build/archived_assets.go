package build

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (b *Builder) buildExtractedAssets() error {
	err := b.buildAssetsArchive()
	if err != nil {
		return errors.WithStack(err)
	}
	return b.disableAssetsHandling()
}

func (b *Builder) disableAssetsHandling() error {
	logrus.Debug("disable asset handling in binary")

	b.transform("actions/app.go", func(body []byte, w io.Writer) error {
		body = bytes.Replace(body, []byte("app.ServeFiles(\"/assets\""), []byte("//app.ServeFiles(\"/assets\""), 1)
		_, err := w.Write(body)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})

	return nil
}

func (b *Builder) buildAssetsArchive() error {
	outputDir := filepath.Dir(filepath.Join(b.Root, b.Bin))

	target := filepath.Join(outputDir, "assets.zip")
	source := filepath.Join(b.Root, "public", "assets")

	logrus.Debugf("building assets archive to %s", target)

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
