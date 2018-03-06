package build

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func (b *Builder) prepTarget() error {
	// Create output directory if not exists
	outputDir := filepath.Join(b.Root, filepath.Dir(b.Bin))
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.MkdirAll(outputDir, 0776)
		logrus.Debugf("creating target dir %s", outputDir)
	}
	return nil
}
