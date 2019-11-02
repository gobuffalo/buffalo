package binding_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/buffalo/render"
	"github.com/stretchr/testify/require"
)

type WithFile struct {
	MyFile binding.File
}

type NamedFile struct {
	MyFile binding.File `form:"afile"`
}

func App() *buffalo.App {
	a := buffalo.New(buffalo.Options{})
	a.POST("/on-struct", func(c buffalo.Context) error {
		wf := &WithFile{}
		if err := c.Bind(wf); err != nil {
			return err
		}
		return c.Render(http.StatusCreated, render.String(wf.MyFile.Filename))
	})
	a.POST("/named-file", func(c buffalo.Context) error {
		wf := &NamedFile{}
		if err := c.Bind(wf); err != nil {
			return err
		}
		return c.Render(http.StatusCreated, render.String(wf.MyFile.Filename))
	})
	a.POST("/on-context", func(c buffalo.Context) error {
		f, err := c.File("MyFile")
		if err != nil {
			return err
		}
		return c.Render(http.StatusCreated, render.String(f.Filename))
	})

	return a
}

func Test_File_Upload_On_Struct(t *testing.T) {
	r := require.New(t)

	req, err := newfileUploadRequest("/on-struct", "MyFile", "file_test.go")
	r.NoError(err)
	res := httptest.NewRecorder()

	App().ServeHTTP(res, req)

	r.Equal(http.StatusCreated, res.Code)
	r.Equal("file_test.go", res.Body.String())
}

func Test_File_Upload_On_Struct_WithTag(t *testing.T) {
	r := require.New(t)

	req, err := newfileUploadRequest("/named-file", "afile", "file_test.go")
	r.NoError(err)
	res := httptest.NewRecorder()

	App().ServeHTTP(res, req)

	r.Equal(http.StatusCreated, res.Code)
	r.Equal("file_test.go", res.Body.String())
}

func Test_File_Upload_On_Context(t *testing.T) {
	r := require.New(t)

	req, err := newfileUploadRequest("/on-context", "MyFile", "file_test.go")
	r.NoError(err)
	res := httptest.NewRecorder()

	App().ServeHTTP(res, req)

	r.Equal(http.StatusCreated, res.Code)
	r.Equal("file_test.go", res.Body.String())
}

// this helper method was inspired by this blog post by Matt Aimonetti:
// https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
func newfileUploadRequest(uri string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
