package alchemist

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TestAdeptMakeDir tests the directory creation in testdata.
func TestAdeptMakeDir(t *testing.T) {

	helper := newAdept(log.New(io.Discard, "", 0), Options{})

	testCases := []struct {
		name    string
		baseDir string
		wantErr error
	}{{
		name:    "ok",
		baseDir: filepath.Join(TestDataDir, "dir"),
	}, {
		name:    "missing write access",
		baseDir: filepath.Join(noAccessDir, "dir"),
		wantErr: IOError{
			Cmd: "make dir",
			Arg: filepath.Join(noAccessDir, "dir", "subdir"),
		},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			dir := filepath.Join(c.baseDir, "subdir")
			err := os.RemoveAll(c.baseDir)
			if err != nil {
				t.Fatalf("ERROR: test setup failed: %v", err)
			}
			if !testing.Verbose() {
				defer os.RemoveAll(c.baseDir)
			}

			err = helper.makedir(dir)
			check.Error(t, err, c.wantErr,
				cmpopts.IgnoreFields(IOError{}, "Err"))

			fileInfo, err := os.Stat(dir)
			if err != nil {
				t.Fatalf("ERROR: got error: %v", err)
			}
			if !fileInfo.IsDir() {
				t.Errorf("ERROR: not a directory")
			}
		})
	}
}

// TestAdeptCopyFileError tests the error behavior of copyFile.
// The behavior of copy is tested in TestAdeptCopyMove.
func TestAdeptCopyFileError(t *testing.T) {

	helper := newAdept(log.New(io.Discard, "", 0), Options{})

	testCases := []struct {
		name     string
		from, to string
		wantErr  error
	}{{
		name:    "source not found",
		from:    notExistName,
		wantErr: IOError{Cmd: "open", Arg: notExistName},
	}, {
		name: "target not usable",
		from: filepath.Join(TestDataDir, fromName),
		to:   filepath.Join(notExistName, notExistName),
		wantErr: IOError{
			Cmd: "create",
			Arg: filepath.Join(notExistName, notExistName),
		},
	}, {
		name: "copy not possible from directory",
		from: TestDataDir,
		to:   filepath.Join(TestDataDir, toName),
		wantErr: IOError{
			Cmd: "copy",
			Arg: TestDataDir + " - " + filepath.Join(TestDataDir, toName),
		},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := helper.copyFile(c.from, c.to)
			check.Error(t, err, c.wantErr,
				cmpopts.IgnoreFields(IOError{}, "Err"))
		})
	}
}

func TestExamine(t *testing.T) {

	testCases := []struct {
		name          string
		from, to      string
		want, wantDir string
		wantRec       bool
		wantErr       error
	}{{
		name: "file to file",
		from: filepath.Join(TestDataDir, fromName),
		to:   filepath.Join(TestDataDir, "workflow.yaml"),
		want: filepath.Join(TestDataDir, "workflow.yaml"),
	}, {
		name: "file to existing dir",
		from: filepath.Join(TestDataDir, fromName),
		to:   filepath.Join(TestDataDir, "page1", "dir1"),
		want: filepath.Join(TestDataDir, "page1", "dir1", fromName),
	}, {
		name:    "file to not existing file",
		from:    filepath.Join(TestDataDir, fromName),
		to:      filepath.Join(TestDataDir, notExistName),
		want:    filepath.Join(TestDataDir, notExistName),
		wantDir: TestDataDir,
	}, {
		name: "file to not existing dir",
		from: filepath.Join(TestDataDir, fromName),
		to: filepath.Join(TestDataDir, "page1", notExistName) +
			string(os.PathSeparator),
		want:    filepath.Join(TestDataDir, "page1", notExistName, fromName),
		wantDir: filepath.Join(TestDataDir, "page1", notExistName),
	}, {
		name:    "dir to existing dir",
		from:    filepath.Join(TestDataDir, "page1"),
		to:      filepath.Join(TestDataDir, "page2"),
		want:    filepath.Join(TestDataDir, "page2", "page1"),
		wantDir: filepath.Join(TestDataDir, "page2", "page1"),
		wantRec: true,
	}, {
		name: "dir to not existing dir",
		from: filepath.Join(TestDataDir, "page1"),
		to: filepath.Join(TestDataDir, "page2", notExistName) +
			string(os.PathSeparator),
		want:    filepath.Join(TestDataDir, "page2", notExistName),
		wantDir: filepath.Join(TestDataDir, "page2", notExistName),
		wantRec: true,
	}, {
		name: "source not found",
		from: filepath.Join(TestDataDir, notExistName),
		to:   filepath.Join(TestDataDir, "workflow.yaml"),
		wantErr: IOError{
			Cmd: "stat",
			Arg: filepath.Join(TestDataDir, notExistName),
		},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			got, gotDir, gotRec, err := examine(c.from, c.to)
			check.Error(t, err, c.wantErr,
				cmpopts.IgnoreFields(IOError{}, "Err"))

			if got != c.want {
				t.Errorf("ERROR: target: got %v, want %v", got, c.want)
			}
			if gotDir != c.wantDir {
				t.Errorf("ERROR: dir: got %v, want %v", gotDir, c.wantDir)
			}
			if gotRec != c.wantRec {
				t.Errorf("ERROR: recursive: got %v, want %v", gotRec, c.wantRec)
			}
		})
	}
}

const (
	srcFileContent = "hello 世界"
	fromName       = "source.txt"
	toName         = "target.txt"
	notExistName   = "does not exist.txt"
)

// TestAdeptFileCopy tests copying and moving a single file.
func TestAdeptFileCopy(t *testing.T) {

	testCases := []struct {
		name     string
		from, to string
		want     string
		wantErr  error
	}{{
		name: "copy ok",
		from: filepath.Join(TestDataDir, fromName),
		to:   filepath.Join(TestDataDir, toName),
		want: filepath.Join(TestDataDir, toName),
	}, {
		name: "copy target is directory",
		from: filepath.Join(TestDataDir, fromName),
		to:   filepath.Join(TestDataDir, "page1"),
		want: filepath.Join(TestDataDir, "page1", fromName),
	}, {
		name:    "copy source does not exist",
		from:    notExistName,
		to:      filepath.Join(TestDataDir, toName),
		wantErr: IOError{Cmd: "stat", Arg: notExistName},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			validateFileCopy(t, c)
		})
	}
}

// validateFileCopy executes the file copy test.
// It is also called by the docker tests.
func validateFileCopy(t *testing.T, c struct {
	name    string
	from    string
	to      string
	want    string
	wantErr error
}) {
	t.Helper()

	err := os.Remove(c.want)
	if err != nil {
		// no error here, file does not exist on the first run
		// or is not writeable
		t.Logf("INFO: test preparation: %v", err)
	}
	if !testing.Verbose() {
		defer os.Remove(c.want)
	}

	helper := newAdept(log.New(io.Discard, "", 0), Options{})

	err = helper.copy(c.from, c.to)
	check.Error(t, err, c.wantErr,
		cmpopts.IgnoreFields(IOError{}, "Err"))

	got, err := os.ReadFile(c.want)
	if err != nil {
		t.Errorf("ERROR: got error: %v", err)
	}

	if diff := cmp.Diff(strings.TrimSpace(string(got)), srcFileContent); diff != "" {
		t.Errorf("ERROR: got- want+\n%s\n", diff)
	}
}

func TestAdeptDirCopy(t *testing.T) {

	newPage := "new_page"
	existingPage := "existing_page"

	helper := newAdept(log.New(io.Discard, "", 0), Options{})

	wantedFiles := []string{
		FormulaFileName,
		"file1.txt",
		"file2.txt",
		filepath.Join("dir1", "file3.txt"),
	}

	testCases := []struct {
		name     string
		from, to string
		prepare  func() error
		wantDir  string
	}{
		{
			name: "new dir",
			from: filepath.Join(TestDataDir, "page1"),
			to:   filepath.Join(TestDataDir, newPage),
			prepare: func() error {
				os.RemoveAll(filepath.Join(TestDataDir, newPage))
				return nil
			},
			wantDir: filepath.Join(TestDataDir, newPage),
		},
		{
			name: "existing dir",
			from: filepath.Join(TestDataDir, "page1"),
			to:   filepath.Join(TestDataDir, existingPage),
			prepare: func() error {
				os.RemoveAll(filepath.Join(TestDataDir, existingPage))
				return os.Mkdir(filepath.Join(TestDataDir, existingPage), dirMode)
			},
			wantDir: filepath.Join(TestDataDir, existingPage, "page1"),
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			err := c.prepare()
			if err != nil {
				t.Fatalf("ERROR: test preparation failed")
			}
			defer os.RemoveAll(c.to)

			err = helper.copy(c.from, c.to)
			if err != nil {
				t.Errorf("ERROR: got error: %v", err)
			}
			if c.wantDir != "" {
				for i, file := range wantedFiles {
					if _, err = os.Stat(filepath.Join(c.wantDir, file)); err != nil {
						t.Errorf("ERROR: file %d not found: %v", i+1, err)
					}
				}
			}
		})
	}
}
