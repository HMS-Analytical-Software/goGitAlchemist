package alchemist_test

// This file contains a black box test.
// It exists in the package alchemist_test which is the special
// black box test package for package alchemist.
// It has no privileged access to package alchemist so it has to import it.

import (
	"path/filepath"
	"testing"

	"github.com/HMS-Analytical-Software/goGitAlchemist/pkg/alchemist"
	"github.com/HMS-Analytical-Software/goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestListPages(t *testing.T) {

	testCases := []struct {
		name     string
		taskList []string
		want     []string
		wantErr  error
	}{{
		name:     "one file",
		taskList: []string{"page1"},
		want: []string{
			filepath.Join(alchemist.TestDataDir, "page1", alchemist.FormulaFileName),
		},
	}, {
		name:     "two files",
		taskList: []string{"page1", "page2"},
		want: []string{
			filepath.Join(alchemist.TestDataDir, "page1", alchemist.FormulaFileName),
			filepath.Join(alchemist.TestDataDir, "page2", alchemist.FormulaFileName),
		},
	}, {
		name:     "file not found",
		taskList: []string{"page1", "page3"},
		wantErr: alchemist.IOError{
			Cmd: "stat",
			Arg: filepath.Join(alchemist.TestDataDir, "page3", alchemist.FormulaFileName),
		},
	}, {
		name:     "dir not found",
		taskList: []string{"page1", "notexist", "page3"},
		wantErr: alchemist.IOError{
			Cmd: "stat",
			Arg: filepath.Join(alchemist.TestDataDir, "notexist", alchemist.FormulaFileName),
		},
	}, {
		name:     "empty list",
		taskList: []string{},
		want:     []string{},
	}, {
		name:     "nil list",
		taskList: []string{},
		want:     []string{},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			got, err := alchemist.ListPages(alchemist.TestDataDir, c.taskList...)

			check.Error(t, err, c.wantErr,
				cmpopts.IgnoreFields(alchemist.IOError{}, "Err"))

			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("ERROR: got- want+\n%s\n", diff)
			}
		})
	}
}

func TestListBookContent(t *testing.T) {

	testCases := []struct {
		name    string
		dir     string
		want    []string
		wantErr error
	}{{
		name: "ok",
		dir:  alchemist.TestDataDir,
		want: []string{
			filepath.Join(alchemist.TestDataDir, "page1", alchemist.FormulaFileName),
			filepath.Join(alchemist.TestDataDir, "page2", alchemist.FormulaFileName),
		},
	}, {
		name: "no subdirs",
		dir:  filepath.Join(alchemist.TestDataDir, "page1"),
	}, {
		name:    "dir not found",
		dir:     "not found",
		wantErr: alchemist.IOError{Cmd: "read dir", Arg: "not found"},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			got, err := alchemist.ListBookContent(c.dir)
			check.Error(t, err, c.wantErr,
				cmpopts.IgnoreFields(alchemist.IOError{}, "Err"))

			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("ERROR: got- want+\n%s\n", diff)
			}
		})
	}
}
