package check

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// ErrorString compares a received error with an expected error.
//
// If want is not nil:
//   - it is checked if the received error is not nil.
//   - the errors are compared, using the prodvide cmp options
//   - the rest of the test is skipped.
//
// If want is nil, the function checks if err is non-nil.
func Error(t errStringTester, err, want error, opt ...cmp.Option) {
	t.Helper()

	if want == nil {
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

	if diff := cmp.Diff(err, want, opt...); diff != "" {
		t.Fatalf("ERROR: got- want+\n%s\n", diff)
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
