package vfsutil_test

import (
	"fmt"
	"log"
	"os"

	"github.com/shurcooL/go/vfs/godocfs/vfsutil"
	"golang.org/x/tools/godoc/vfs"
	"golang.org/x/tools/godoc/vfs/mapfs"
)

func ExampleWalk() {
	var fs vfs.FileSystem = mapfs.New(map[string]string{
		"zzz-last-file.txt":   "It should be visited last.",
		"a-file.txt":          "It has stuff.",
		"another-file.txt":    "Also stuff.",
		"folderA/entry-A.txt": "Alpha.",
		"folderA/entry-B.txt": "Beta.",
	})

	walkFn := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Printf("can't stat file %s: %v\n", path, err)
			return nil
		}
		fmt.Println(path)
		return nil
	}

	err := vfsutil.Walk(fs, "/", walkFn)
	if err != nil {
		panic(err)
	}

	// Output:
	// /
	// /a-file.txt
	// /another-file.txt
	// /folderA
	// /folderA/entry-A.txt
	// /folderA/entry-B.txt
	// /zzz-last-file.txt
}
