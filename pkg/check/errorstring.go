// Package check provides unit test helper code.
package check

import (
	"strings"
	"testing"
)

// ErrorString compares a received error with an expected error message.
//
// If want is not the empty string:
//   - it is checked if the error is not nil.
//   - it is checked if the error message contains the wanted string.
//   - the rest of the test is skipped.
//
// If want is empty, the function checks if err is non-nil.
func ErrorString(t errStringTester, err error, want string) {
	t.Helper()

	if want == "" {
		// no error expected.
		if err != nil {
			t.Fatalf("ERROR: got error: %v", err)
		}
		return
	}

	if err == nil {
		// expected error not detected.
		t.Fatalf("ERROR: error not detected")
		return
	}

	if !strings.Contains(err.Error(), want) {
		// error message does not match the expectation.
		t.Fatalf("ERROR: got %s, want %s", err.Error(), want)
		return
	}

	// test passed.
	if testing.Verbose() {
		t.Logf("INFO: got error: %v", err)
	}

	// if we expect an error, the test statements following
	// the check.ErrorString call should not be executed.
	t.Skip("skip subsequent statements due to wanted error")
}

// errStringTester defined the subset of *testing.T methods used here.
type errStringTester interface {
	Helper()
	Skip(...any)
	Fatalf(string, ...any)
	Logf(string, ...any)
}
