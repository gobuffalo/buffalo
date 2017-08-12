package builder

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type box struct {
	Name  string
	Files []file
}

func (b *box) Walk(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		name := strings.Replace(path, root+string(os.PathSeparator), "", 1)
		name = strings.Replace(name, "\\", "/", -1)
		f := file{
			Name: name,
		}

		bb, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		bb, err = json.Marshal(bb)
		if err != nil {
			return errors.WithStack(err)
		}
		f.Contents = strings.Replace(string(bb), "\"", "\\\"", -1)

		b.Files = append(b.Files, f)
		return nil
	})
}
