//go:build acctest

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/HMS-Analytical-Software/goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
)

const testDataDir = "testdata"

// TestAcceptanceRunAll tests the -runall flag.
//
// It uses the same test case definition as TestAcceptanceTask.
func TestAcceptanceRunAll(t *testing.T) {

	// build gitalchemist binary to test
	err := exec.Command(goCmd, "build").Run()
	if err != nil {
		// should not happen, as go test builds the TEST binary first.
		t.Fatalf("ERROR: test preparation: %v", err)
	}

	// call gitalchemist only once.
	err = callGitAlchemist(t, "-runall")
	if err != nil {
		t.Fatalf("ERROR: got error %v", err)
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			if c.wantErr != "" {
				// ignore error test cases here
				return
			}

			checkFormulaResult(t, c)
		})
	}
}

// TestAcceptanceTask tests single formula (task) execution.
//
// It uses the same test case definition as TestAcceptanceRunAll.
func TestAcceptanceTask(t *testing.T) {

	// build gitalchemist binary to test
	err := exec.Command(goCmd, "build").Run()
	if err != nil {
		// should not happen, as go test builds the TEST binary first.
		t.Fatalf("ERROR: test preparation: %v", err)
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			// call gitalchemist
			err := callGitAlchemist(t, c.name)
			check.ErrorString(t, err, c.wantErr)

			checkFormulaResult(t, c)
		})
	}
}

// checkFormulaResult executes the tests for a single formula result.
func checkFormulaResult(t *testing.T, c accTestCase) {
	gitDir := filepath.Join(defaultCwd, c.name)
	srcDir := filepath.Join(testDataDir, c.name)

	// bare repo must always be created
	checkGit(t, filepath.Join(defaultCwd, "remotes", c.name),
		"true\n", "rev-parse", "--is-bare-repository")

	// remote repo must always be set
	checkGit(t, gitDir,
		"origin\t"+filepath.Join("..", "remotes", c.name)+" (fetch)\n"+
			"origin\t"+filepath.Join("..", "remotes", c.name)+" (push)\n",
		"remote", "-v")

	// compare file content
	for _, pair := range c.compareList {
		compareFiles(t,
			filepath.Join(srcDir, pair.from),
			filepath.Join(gitDir, pair.to),
		)
	}

	// check git info
	for _, g := range c.gitList {
		checkGit(t, gitDir, g.want, g.args...)
	}
}

// accTestCase defines a single acceptance test case.
//
// It contains
//   - a test name that corresponds to a directory within testdata (= a task)
//   - a list of file pairs that should have the same content
//   - a list of git commands that have to return the expected output
//   - an error message (empty if no error should happen)
type accTestCase struct {
	name        string     // test == task name
	compareList []filePara // files to compare
	gitList     []gitPara  // git test calls
	wantErr     string     // expected error
}

// gitPara contains parameter for the git test calls
type gitPara struct {
	args []string // git call
	want string   // expected result
}

// filePara  contains the two file names that should be compared.
type filePara struct {
	from string // file name in config directory
	to   string // file name in target (git) directory
}

// testCases are used for the single-call test and for the runall-test.
//
// The test cases are extracted from the gitalchemist.yaml file in the
// corresponding testdata directory. The file defines the steps to execute,
// so they can be used to extract what has to be tested.
//
// TODO: the git test calls may be extended and/or optimized.
// Currently, they only test each step in a basic way, but not every property
// that the resulting git repo should have.
var testCases = []accTestCase{
	{
		name: "basic_workflow",
		compareList: []filePara{{
			from: filepath.Join("files", "project_plan_v1.md"),
			to:   "project_plan.md",
		}},
		gitList: []gitPara{{
			// check last file additions
			args: []string{"show", "--pretty=", "--name-status"},
			want: "A\tproject_plan.md\n",
		}, {
			// check log messages
			args: []string{"log", "--pretty=format:%s"},
			want: "Added first file",
		}},
	},
	{
		name: "cmd_create_add_commit",
		compareList: []filePara{{
			from: filepath.Join("files", "project_plan_v3.md"),
			to:   "project_plan.md",
		}, {
			from: filepath.Join("files", "folder1", "file1.md"),
			to:   filepath.Join("folder1", "file1.md"),
		}, {
			from: filepath.Join("files", "folder1", "file2.md"),
			to:   filepath.Join("folder1", "file2.md"),
		}, {
			from: filepath.Join("files", "folder1", "folder2", "file3.md"),
			to:   filepath.Join("folder1", "folder2", "file3.md"),
		}},
		gitList: []gitPara{{
			args: []string{"show", "--pretty=", "--name-status"},
			want: "A\tfolder1/file1.md\n" +
				"A\tfolder1/file2.md\n" +
				"A\tfolder1/folder2/file3.md\n",
		}, {
			args: []string{"log", "--pretty=format:%s"},
			want: "added folder1\n" +
				"removed unnecessary parts of the project plan\n" +
				"added project summary\n" +
				"added project plan template",
		}, {
			args: []string{"status"}, // ensure push
			want: "On branch main\n" +
				"Your branch is up to date with 'origin/main'.\n\n" +
				"nothing to commit, working tree clean\n",
		}},
	},
	{
		name: "cmd_create_file",
		compareList: []filePara{{
			from: filepath.Join("files", "project_plan_v1.md"),
			to:   "project_plan.md",
		}, {
			from: filepath.Join("files", "some_other_file.txt"),
			to:   "some_other_file.txt",
		}},
	},
	{
		name: "cmd_merge",
		compareList: []filePara{{
			from: filepath.Join("files", "main_v1.py"),
			to:   "main.py",
		}, {
			from: filepath.Join("files", "readme_v1.md"),
			to:   "readme.md",
		}, {
			from: filepath.Join("files", "gitignore_file"),
			to:   ".gitignore",
		}},
		gitList: []gitPara{{
			args: []string{"show", "--pretty=", "--name-status"},
			want: "A\t.gitignore\nA\tmain.py\n",
		}, {
			args: []string{"log", "--pretty=format:%s"},
			want: "initial commit password generator project\n" +
				"readme",
		}, {
			args: []string{"branch"}, // check branches
			want: "  feature/start_project\n" +
				"* main\n",
			// TODO: ensure merge
		}},
	},
	{
		name: "cmd_mv",
		compareList: []filePara{{
			from: filepath.Join("files", "readme_v1.md"),
			to:   "readme.md",
		}, {
			from: filepath.Join("files", "main_v1.py"),
			to:   "generator.py",
		}, {
			from: filepath.Join("files", "gitignore_file"),
			to:   ".gitignore",
		}},
	},
	{
		name: "cmd_remove_and_commit",
		compareList: []filePara{{
			from: filepath.Join("files", "hello.py"),
			to:   "hello.py",
		}},
		gitList: []gitPara{{
			args: []string{"show", "--pretty=", "--name-status"},
			want: "D\tnotes-timeline.txt\n",
		}, {
			args: []string{"log", "--pretty=format:%s"},
			want: "clean up timeline notes\nhello world",
		}},
	},
	{
		name:    "error case this does not exist",
		wantErr: "exit status 6",
	},
}

// compareFiles compares two files and emits an error message if they dont match.
//
// The files should not be big, because the content is read completely at once.
func compareFiles(t *testing.T, file1, file2 string) {
	t.Helper()

	content1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("ERROR: read file: %v", err)
	}
	content2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("ERROR: read file: %v", err)
	}

	if string(content1) != string(content2) {
		t.Fatalf("ERROR: files differ: %q - %q", file1, file2)
	}

	t.Logf("INFO: files equal: %q - %q", file1, file2)
}

// checkGit calls git with the provided parameters in the provided
// directory. It checks its exit status and the expected output.
//
// In Verbose mode, stdout and stderr are displayed.
func checkGit(t *testing.T, dir, want string, args ...string) {
	t.Helper()

	cmd := exec.Command(gitCmd, args...)
	cmd.Dir = dir
	got, err := cmd.CombinedOutput()
	t.Logf("INFO: %v", args)
	t.Logf("INFO: %s", got)

	if err != nil {
		t.Errorf("ERROR: got error: %v", err)
	}
	if diff := cmp.Diff(string(got), want); diff != "" {
		t.Errorf("ERROR: got- want+\n%s\n", diff)
	}
}

// callGitAlchemist  calls gitalchemist with the provided parameter and
// returns its exit status.
//
// parameter may contain:
//   - a task name
//   - parameter -runall
//
// In Verbose mode, stdout and stderr are displayed.
func callGitAlchemist(t *testing.T, parameter string) error {
	t.Helper()

	cmd := exec.Command(gitAlchemistCmd, "-cfgdir", testDataDir, parameter)
	if testing.Verbose() {
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
