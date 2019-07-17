package render

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Template(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	box := e.TemplatesBox
	r.NoError(box.AddString(htmlTemplate, `<%= name %>`))

	re := e.Template("foo/bar", htmlTemplate)
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Mark"}))
}

//
// func Test_AssetPath(t *testing.T) {
// 	r := require.New(t)
//
// 	cases := map[string]string{
// 		"something.txt":         "/assets/something.txt",
// 		"images/something.png":  "/assets/images/something.png",
// 		"/images/something.png": "/assets/images/something.png",
// 		"application.css":       "/assets/application.aabbc123.css",
// 	}
//
// 	tDir, err := ioutil.TempDir("", "templates")
// 	if err != nil {
// 		r.Fail("Could not set the templates dir")
// 	}
//
// 	aDir, err := ioutil.TempDir("", "assets")
// 	if err != nil {
// 		r.Fail("Could not set the assets dir")
// 	}
//
// 	re := render.New(render.Options{
// 		TemplatesBox: packr.New(tDir, tDir),
// 		AssetsBox:    packr.New(aDir, aDir),
// 	}).Template
//
// 	ioutil.WriteFile(filepath.Join(aDir, "manifest.json"), []byte(`{
// 		"application.css": "application.aabbc123.css"
// 	}`), 0644)
//
// 	for original, expected := range cases {
//
// 		tmpFile, err := os.Create(filepath.Join(tDir, "test.html"))
// 		r.NoError(err)
//
// 		_, err = tmpFile.Write([]byte("<%= assetPath(\"" + original + "\") %>"))
// 		r.NoError(err)
//
// 		result := re("text/html; charset=utf-8", filepath.Base(tmpFile.Name()))
//
// 		bb := &bytes.Buffer{}
// 		err = result.Render(bb, render.Data{})
// 		r.NoError(err)
// 		r.Equal(expected, strings.TrimSpace(bb.String()))
//
// 		os.Remove(tmpFile.Name())
// 	}
// }
//
// func Test_AssetPathNoManifest(t *testing.T) {
// 	r := require.New(t)
//
// 	cases := map[string]string{
// 		"something.txt": "/assets/something.txt",
// 	}
//
// 	tDir, err := ioutil.TempDir("", "templates")
// 	if err != nil {
// 		r.Fail("Could not set the templates dir")
// 	}
//
// 	aDir, err := ioutil.TempDir("", "assets")
// 	if err != nil {
// 		r.Fail("Could not set the assets dir")
// 	}
//
// 	re := render.New(render.Options{
// 		TemplatesBox: packr.New(tDir, tDir),
// 		AssetsBox:    packr.New(aDir, aDir),
// 	}).Template
//
// 	for original, expected := range cases {
//
// 		tmpFile, err := os.Create(filepath.Join(tDir, "test.html"))
// 		r.NoError(err)
//
// 		_, err = tmpFile.Write([]byte("<%= assetPath(\"" + original + "\") %>"))
// 		r.NoError(err)
//
// 		result := re("text/html; charset=utf-8", filepath.Base(tmpFile.Name()))
//
// 		bb := &bytes.Buffer{}
// 		err = result.Render(bb, render.Data{})
// 		r.NoError(err)
// 		r.Equal(expected, strings.TrimSpace(bb.String()))
//
// 		os.Remove(tmpFile.Name())
// 	}
// }
// func Test_AssetPathManifestCorrupt(t *testing.T) {
// 	r := require.New(t)
//
// 	cases := map[string]string{
// 		"something.txt": "manifest.json is not correct",
// 		"other.txt":     "manifest.json is not correct",
// 	}
//
// 	tDir, err := ioutil.TempDir("", "templates")
// 	if err != nil {
// 		r.Fail("Could not set the templates dir")
// 	}
//
// 	aDir, err := ioutil.TempDir("", "assets")
// 	if err != nil {
// 		r.Fail("Could not set the assets dir")
// 	}
//
// 	ioutil.WriteFile(filepath.Join(aDir, "manifest.json"), []byte(`//shdnn Corrupt!`), 0644)
//
// 	re := render.New(render.Options{
// 		TemplatesBox: packr.New(tDir, tDir),
// 		AssetsBox:    packr.New(aDir, aDir),
// 	}).Template
//
// 	for original, expected := range cases {
//
// 		tmpFile, err := os.Create(filepath.Join(tDir, "test.html"))
// 		r.NoError(err)
//
// 		_, err = tmpFile.Write([]byte("<%= assetPath(\"" + original + "\") %>"))
// 		r.NoError(err)
//
// 		result := re("text/html; charset=utf-8", filepath.Base(tmpFile.Name()))
//
// 		bb := &bytes.Buffer{}
// 		err = result.Render(bb, render.Data{})
// 		r.Error(err)
// 		r.Contains(err.Error(), expected)
//
// 		os.Remove(tmpFile.Name())
// 	}
// }
