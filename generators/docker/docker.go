package docker

import (
	"github.com/gobuffalo/buffalo/generators"
	"github.com/gobuffalo/makr"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

func dockerFiles(d Generator) (generators.Files, error) {
	var box packr.Box

	switch d.Style {
	case "standard":
		box = packr.NewBox("./standard/templates")
	case "multi":
		box = packr.NewBox("./multi/templates")
	default:
		return nil, errors.Errorf("unknown Docker style: %s", d.Style)
	}

	files, err := generators.FindByBox(box)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return files, nil
}

func dockerComposeFiles(d Generator) (generators.Files, error) {
	var box packr.Box

	switch d.DockerCompose {
	case "deps":
		box = packr.NewBox("./compose-deps")
	case "full":
		box = packr.NewBox("./compose-full")
	default:
		return nil, errors.Errorf("unknown docker-compose style: %s", d.DockerCompose)
	}

	files, err := generators.FindByBox(box)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return files, nil
}

// Run Docker generator
func (d Generator) Run(root string, data makr.Data) error {
	g := makr.New()
	data["opts"] = d
	g.Add(&makr.Func{
		Should: func(data makr.Data) bool {
			return d.Style != "none"
		},
		Runner: func(root string, data makr.Data) error {
			files, err := dockerFiles(d)
			if err != nil {
				return errors.WithStack(err)
			}

			dockerCompose, err := dockerComposeFiles(d)
			if err != nil {
				return errors.WithStack(err)
			}

			files = append(files, dockerCompose...)
			fg := makr.New()
			for _, f := range files {
				fg.Add(makr.NewFile(f.WritePath, f.Body))
			}
			return fg.Run(root, data)
		},
	})
	return g.Run(root, data)
}
