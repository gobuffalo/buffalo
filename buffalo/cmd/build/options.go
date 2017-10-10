package build

import (
	"encoding/json"

	"github.com/gobuffalo/buffalo/meta"
)

type Options struct {
	meta.App
	ExtractAssets bool     `json:"extract_assets"`
	HasDB         bool     `json:"has_db"`
	LDFlags       string   `json:"ld_flags"`
	Tags          []string `json:"tags"`
	Static        bool     `json:"static"`
	Debug         bool     `json:"debug"`
	Compress      bool     `json:"compress"`
}

func (o Options) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}
