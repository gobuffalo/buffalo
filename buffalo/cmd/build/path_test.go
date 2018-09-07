package build

import (
	"testing"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/stretchr/testify/require"
)

const rootPath = "/notabsolutepath"

func TestAbsoluteBinaryPath(t *testing.T) {
	tests := []struct {
		binaryPath           string
		expectedAbsolutePath string
	}{
		// relative paths
		{"binary", rootPath + "/binary"},
		{"something/else", rootPath + "/something/else"},

		// absolute paths
		{"/binary", "/binary"},
		{"/tmp/binary", "/tmp/binary"},
	}

	for _, test := range tests {
		b := Builder{
			Options: Options{
				App: meta.App{
					Root: rootPath,
					Bin:  test.binaryPath,
				},
			},
		}
		require.Equal(t, test.expectedAbsolutePath, b.AbsoluteBinaryPath())
	}
}
