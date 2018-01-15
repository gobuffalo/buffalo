package build

import (
	"encoding/json"

	"github.com/gobuffalo/buffalo/meta"
)

// Options for a build
type Options struct {
	meta.App
	ExtractAssets          bool     `json:"extract_assets"`
	WithAssets             bool     `json:"with_assets"`
	LDFlags                string   `json:"ld_flags"`
	Tags                   []string `json:"tags"`
	Static                 bool     `json:"static"`
	Debug                  bool     `json:"debug"`
	Compress               bool     `json:"compress"`
	Environment            string   `json:"environment"`
	SkipTemplateValidation bool     `json:"skip_template_validation"`
}

func (o Options) String() string {
	b, _ := json.Marshal(o)
	return string(b)
}
