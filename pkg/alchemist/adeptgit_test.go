package alchemist

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"testing"

	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TestAdeptGit tests the error behaviour of the adept git method.
//
// This test executes the go test binary instead of the real git
// command to test the error behaviour of the git method.
//
// It consists of two functions
//   - TestAdeptGit:
//     contains the unit test code
//   - TestCalledByAdeptGit:
//     is called by TestAdeptGit. It exits with an error code if it was
//     called with the command line parameter "ERROR",
//     else it exits with status ok.
func TestAdeptGit(t *testing.T) {

	os.Setenv(envFlag, "1")
	testExe, err := os.Executable()
	if err != nil {
		t.Fatalf("ERROR: test setup failed: %v", err)
	}

	// as we use the go test binary, we have to get rid of the
	// additional output.
	goTestOutput := regexp.MustCompile(`(?s)\[DEBUG\] PASS.*`)

	testCases := []struct {
		name    string
		dir     string
		verbose bool
		call    []string
		want    string
		wantErr error
	}{{
		name: "ok",
		call: []string{testFlag},
	}, {
		name:    "verbose",
		verbose: true,
		call:    []string{testFlag},
		want: `[DEBUG] "": git []string{"` + testFlag + `"}
[DEBUG] ok
`,
	}, {
		name:    "dir verbose",
		verbose: true,
		dir:     TestDataDir,
		call:    []string{testFlag},
		want: `[DEBUG] "` + TestDataDir + `": git []string{"` +
			testFlag + `"}
[DEBUG] ok
`,
	}, {
		name:    "with error",
		call:    []string{testFlag, "ERROR"},
		wantErr: ExecError{Cmd: gitCmd, Args: []string{testFlag, "ERROR"}},
	}, {
		name:    "error verbose",
		verbose: true,
		call:    []string{testFlag, "ERROR"},
		wantErr: ExecError{Cmd: gitCmd, Args: []string{testFlag, "ERROR"}},
		want: `[DEBUG] "": git []string{"` + testFlag + `", "ERROR"}
[DEBUG] failed
`,
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			var buf bytes.Buffer
			helper := newAdept(log.New(&buf, "", 0), Options{Verbose: c.verbose})
			helper.exe = testExe
			err := helper.git(c.dir, c.call...)

			got := buf.String()
			got = goTestOutput.ReplaceAllString(got, "")

			if diff := cmp.Diff(got, c.want); diff != "" {
				t.Errorf("ERROR: got- want+\n%s\n", diff)
			}
			check.Error(t, err, c.wantErr, cmpopts.IgnoreFields(ExecError{}, "Err"))
		})
	}
}

// TestCalledByAdeptGit is called by the TestAdeptGit test.
// If one of the command line flags is "ERROR", it returns with an error code,
// otherwise it returns success.
//
// If it is called by go test directly, it does nothing.
func TestCalledByAdeptGit(t *testing.T) {

	if os.Getenv(envFlag) != "1" {
		// not called by TestAdeptGit, so do nothing
		return
	}

	if slices.Contains(os.Args, "ERROR") {
		// program requested to return an error
		fmt.Printf("failed\n")
		os.Exit(42)
	}

	// successful program execution
	fmt.Printf("ok\n")
}

const (
	// envFlag is used to signal the function that it was called by the
	// other test function, not directly by go test.
	envFlag = "TEST_CALLED_BY_ADEPT_GIT"
	// testFlag is used to restrict the go test call done by TestAdeptGit
	// to execute only the test helper function.
	testFlag = "-test.run=TestCalledByAdeptGit"
)
