package build

import (
	"runtime"
	"testing"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/stretchr/testify/assert"
)

const rootPath = "ROOT"

func TestAbsoluteBinaryPath(t *testing.T) {
	tests := []struct {
		inBinaryPath                             string
		expectedNixResult, expectedWindowsResult string
	}{
		// relative paths on *nix
		{"relnixbinary",
			rootPath + "/relnixbinary", rootPath + "\\relnixbinary"},
		{"rel/nix/binary",
			rootPath + "/rel/nix/binary", rootPath + "\\rel\\nix\\binary"},

		// relative paths on Windows
		{"relwinbinary",
			rootPath + "/relwinbinary", rootPath + "\\relwinbinary"},
		{"rel\\win\\binary",
			rootPath + "/rel\\win\\binary", rootPath + "\\rel\\win\\binary"},

		// absolute paths on *nix
		{"/absnixbinary",
			"/absnixbinary", rootPath + "\\absnixbinary"},
		{"/abs/nix/binary",
			"/abs/nix/binary", rootPath + "\\abs\\nix\\binary"},

		// absolute paths on Windows
		{"C:\\abswinbinary",
			rootPath + "/C:\\abswinbinary", "C:\\abswinbinary"},
		{"C:\\abs\\win\\binary",
			rootPath + "/C:\\abs\\win\\binary", "C:\\abs\\win\\binary"},
	}

	for _, test := range tests {
		b := Builder{
			Options: Options{
				App: meta.App{
					Root: rootPath,
					Bin:  test.inBinaryPath,
				},
			},
		}
		if runtime.GOOS == "windows" {
			assert.Equal(t, test.expectedWindowsResult, b.AbsoluteBinaryPath())
		} else {
			assert.Equal(t, test.expectedNixResult, b.AbsoluteBinaryPath())

		}

	}
}
