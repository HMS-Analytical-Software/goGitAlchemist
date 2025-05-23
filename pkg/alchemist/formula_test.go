package alchemist

import (
	"bytes"
	"log"
	"path/filepath"
	"testing"

	"github.com/HMS-Analytical-Software/goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TestDataDir defines the location of the unit test files.
const TestDataDir = "testdata"

// TestReadFormulaFileError tests the error behavior of reading
// the gitalchemist yaml file.
func TestReadFormulaFileError(t *testing.T) {

	testCases := []struct {
		name     string
		fileName string
		wantErr  error
	}{{
		name:     "ok",
		fileName: filepath.Join(TestDataDir, "workflow.yaml"),
	}, {
		name:     "file not found",
		fileName: "does not exist",
		wantErr: IOError{
			Cmd: "open",
			Arg: "does not exist",
		},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Read(c.fileName)
			check.Error(t, err, c.wantErr, cmpopts.IgnoreFields(IOError{}, "Err"))
		})
	}
}

// TestReadFormulaContent tests the extraction of the content
// of the gitalchemist yaml file.
func TestReadFormulaContent(t *testing.T) {

	testCases := []struct {
		name     string
		fileName string
		want     Formula
		wantErr  string
	}{{
		name:     "basic",
		fileName: filepath.Join(TestDataDir, "workflow.yaml"),
		want:     completeFormula,
	}, {
		name:     "not yaml",
		fileName: filepath.Join(TestDataDir, "workflow.xml"),
		wantErr:  "yaml: unmarshal errors:\n  line 1: cannot unmarshal",
	}, {
		name:     "incomplete yaml",
		fileName: filepath.Join(TestDataDir, "incomplete.yaml"),
		wantErr:  "validate init_bare_repo (1): value for bare is missing",
	}, {
		name:     "invalid yaml",
		fileName: filepath.Join(TestDataDir, "invalid.yaml"),
		wantErr:  "yaml node too small: &yaml.Node{",
	}, {
		name:     "invalid node yaml",
		fileName: filepath.Join(TestDataDir, "invalidnode.yaml"),
		wantErr: "yaml decode Formula: yaml decode node alchemist.initRepoSpell: " +
			"yaml: unmarshal errors",
	}, {
		name:     "unknown yaml element",
		fileName: filepath.Join(TestDataDir, "unknown.yaml"),
		wantErr:  `unkonwn command "not_known"`,
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			got, err := Read(c.fileName)
			check.ErrorString(t, err, c.wantErr)

			diff := cmp.Diff(got, c.want, cmp.AllowUnexported(symbols{}))
			if diff != "" {
				t.Errorf("ERROR: got-, want+\n%v\n", diff)
			}
		})
	}
}

// TestFormulaTransmute tests the log messages that a Transmute call emits.
func TestFormulaTransmute(t *testing.T) {

	testCases := []struct {
		name    string
		formula Formula
		verbose bool
		nSteps  int
		want    string
		wantErr error
	}{{
		name:    "normal",
		formula: completeFormula,
		verbose: false,
		want:    wantTestCompleteLog,
	}, {
		name:    "verbose",
		formula: completeFormula,
		verbose: true,
		want:    wantTestVerboseLog,
	}, {
		name:    "short",
		formula: completeFormula,
		nSteps:  2,
		want:    wantTestShortLog,
	}, {
		name:    "error",
		formula: errorFormula,
		want:    "",
		wantErr: ExecError{Cmd: "cast error"},
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			opt := Options{
				TaskDir:       completeFormula.Title,
				Test:          true,
				ExecuteSpells: c.nSteps,
				Verbose:       c.verbose,
				taskName:      "testtask",
				RepoDir:       "repodir",
				cloneTo:       "cloneto",
				CfgDir:        "cfgdir",
			}

			var buf bytes.Buffer
			err := Transmute(c.formula, opt, log.New(&buf, "", 0))

			check.Error(t, err, c.wantErr, cmpopts.IgnoreFields(ExecError{}, "Err"))

			got := buf.String()

			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("ERROR: got- want+\n%s\n", diff)
			}
		})
	}
}

// errorFormula is a formula that contains a test double caster
// that returns an error on the cast method.
var errorFormula = Formula{
	Title: "test_workflow",
	Commands: symbols{
		cloneTo: "workflow",
		spells:  []caster{errorSpell{}},
	},
}

// errorSpell is a test double
type errorSpell struct{}

func (errorSpell) validate() error { return nil }

func (errorSpell) cast(assistant, Options) error {
	return ExecError{Cmd: "cast error"}

}

// completeFormula corresponds to the content of the file testdata/workflow.yaml.
var completeFormula = Formula{
	Title: "test_workflow",
	Commands: symbols{
		cloneTo: "workflow",
		spells: []caster{
			initRepoSpell{
				Bare:    "remotes/create_add_commit",
				CloneTo: "workflow",
			},
			createFileSpell{
				Source: "files/project_plan_v1.md",
				Target: "project_plan.md",
			},
			addSpell{
				Files: []string{"project_plan.md"},
			},
			commitSpell{
				Message: "Added first file",
				Author:  "red",
			},
			createAddCommitSpell{
				Files:   []string{"files/project_plan_v3.md => project_plan.md"},
				Message: "removed unnecessary parts of the project plan",
				Author:  "red",
			},
			createAddCommitSpell{
				Files:   []string{"files/folder1 => folder1/"},
				Message: "added folder1",
				Author:  "red",
			},
			gitSpell{
				Command: `git commit -m "my message"`,
			},
			gitSpell{
				Command: "push origin main",
			},
			mergeSpell{
				Source: "feature/start_project",
				Target: "main",
			},
			mergeSpell{
				Source:       "feature/other",
				Target:       "main",
				DeleteSource: true,
			},
			pushSpell{
				Main: true,
			},
			moveSpell{
				Source: "main.py",
				Target: "generator.py",
			},
			removeAndCommitSpell{
				Files:   []string{"notes-timeline.txt"},
				Message: "clean up timeline notes",
				Author:  "red",
			},
			removeAndCommitSpell{
				Files:   []string{"notes-timeline_nonexistent.txt"},
				Message: "clean up timeline notes that don't exist (this step should fail)",
				Author:  "red",
			},
		},
	},
}
