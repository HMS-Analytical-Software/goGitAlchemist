//go:build testrundocker

package alchemist

import (
	"path/filepath"
	"testing"
)

// TestDockerAdeptFileCopy executes a file copy test in a docker container.
func TestDockerAdeptFileCopy(t *testing.T) {
	testCases := []struct {
		name     string
		from, to string
		want     string
		wantErr  error
	}{{
		name: "copy target not writeable",
		from: filepath.Join(TestDataDir, fromName),
		to:   filepath.Join("/", toName),
		wantErr: IOError{
			Cmd: "create",
			Arg: filepath.Join("/", toName),
		},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			validateFileCopy(t, c)
		})
	}

}
