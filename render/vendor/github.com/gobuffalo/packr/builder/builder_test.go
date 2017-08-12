package builder

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_Builder_Run(t *testing.T) {
	r := require.New(t)

	root := filepath.Join("..", "example")
	defer Clean(root)

	exPackr := filepath.Join(root, "example-packr.go")
	r.False(fileExists(exPackr))

	fooPackr := filepath.Join(root, "foo", "foo-packr.go")
	r.False(fileExists(fooPackr))

	b := New(context.Background(), root)
	err := b.Run()
	r.NoError(err)

	r.True(fileExists(exPackr))
	r.True(fileExists(fooPackr))

	bb, err := ioutil.ReadFile(exPackr)
	r.NoError(err)
	r.True(bytes.Contains(bb, []byte(`packr.PackJSONBytes("./assets", "app.css", "\"Ym9keSB7CiAgYmFja2dyb3VuZDogcmVkOwp9Cg==\"")`)))
	r.True(bytes.Contains(bb, []byte(`packr.PackJSONBytes("./assets", "app.js", "\"YWxlcnQoImhlbGxvISIpOwo=\"")`)))

	bb, err = ioutil.ReadFile(fooPackr)
	r.NoError(err)
	r.True(bytes.Contains(bb, []byte(`packr.PackJSONBytes("../assets", "app.css", "\"Ym9keSB7CiAgYmFja2dyb3VuZDogcmVkOwp9Cg==\"")`)))
	r.True(bytes.Contains(bb, []byte(`packr.PackJSONBytes("../assets", "app.js", "\"YWxlcnQoImhlbGxvISIpOwo=\"")`)))
	r.True(bytes.Contains(bb, []byte(`packr.PackJSONBytes("../templates", "index.html", "\"PCFET0NUWVBFIGh0bWw+CjxodG1sPgogIDxoZWFkPgogICAgPG1ldGEgY2hhcnNldD0idXRmLTgiIC8+CiAgICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoIiAvPgogICAgPHRpdGxlPklOREVYPC90aXRsZT4KICAgIGxpbmsKICA8L2hlYWQ+CiAgPGJvZHk+CiAgICBib2R5CiAgPC9ib2R5Pgo8L2h0bWw+Cg==\"")`)))
}

func Test_Binary_Builds(t *testing.T) {
	r := require.New(t)
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)

	root := "../example"
	defer Clean(root)
	defer os.RemoveAll(filepath.Join(root, "bin"))

	b := New(context.Background(), root)
	err := b.Run()
	r.NoError(err)

	os.Chdir(root)
	cmd := exec.Command(envy.Get("GO_BIN", "go"), "build", "-v", "-o", "bin/example")
	err = cmd.Run()
	r.NoError(err)

	r.True(fileExists("bin/example"))
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
