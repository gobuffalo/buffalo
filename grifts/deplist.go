package grifts

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/markbates/deplist"
	"github.com/markbates/grift/grift"
)

func depList() []string {
	list, _ := deplist.List("examples")
	clean := []string{}
	for v := range list {
		if !strings.Contains(v, "gobuffalo/buffalo") {
			clean = append(clean, v)
		}
	}
	sort.Strings(clean)
	return clean
}

var _ = grift.Add("deplist", func(c *grift.Context) error {
	w, err := os.Create("deplist")
	if err != nil {
		return err
	}
	defer w.Close()
	w.WriteString(strings.Join(depList(), "\n"))
	return nil
})

var _ = grift.Add("deplist:count", func(c *grift.Context) error {
	fmt.Printf("%d Dependencies\n", len(depList()))
	return nil
})

var _ = grift.Add("deplist:print", func(c *grift.Context) error {
	fmt.Println(strings.Join(depList(), "\n"))
	return nil
})
