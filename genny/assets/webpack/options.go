package webpack

import "github.com/gobuffalo/buffalo/meta"

// Options for creating a new webpack setup
type Options struct {
	meta.App
	Bootstrap int `json:"bootstrap"`
}
