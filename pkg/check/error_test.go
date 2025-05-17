package check_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TestErrorString tests the ErrorString error helper function.
func TestErrorString(t *testing.T) {

	dummyErrMsg := "dummy"
	dummyErr := errors.New(dummyErrMsg)

	testCases := []struct {
		name        string
		gotErr      error
		expectedMsg string
		expectedErr error

		wantErrMsgString string
		wantErrMsg       string
		wantSkipMsg      string
		wantLogMsg       string // set only in verbose test mode
	}{{
		name: "no error", // no messages from ErrorString
	}, {
		name:        "matching error",
		gotErr:      dummyErr,
		expectedMsg: dummyErrMsg,
		expectedErr: dummyErr,
		wantSkipMsg: "[skip subsequent statements due to wanted error]",
		wantLogMsg:  "INFO: got error: " + dummyErrMsg,
	}, {
		name:             "unexpected error",
		gotErr:           dummyErr,
		wantErrMsgString: "ERROR: got error: " + dummyErrMsg,
		wantErrMsg:       "ERROR: got error: " + dummyErrMsg,
	}, {
		name:             "expected error not detected",
		expectedMsg:      dummyErrMsg,
		expectedErr:      dummyErr,
		wantErrMsgString: "ERROR: error not detected",
		wantErrMsg:       "ERROR: error not detected",
	}, {
		name:             "nonmatching error",
		gotErr:           dummyErr,
		expectedMsg:      "does not match",
		expectedErr:      errors.New("does not mnatch"),
		wantErrMsgString: "ERROR: got " + dummyErrMsg + ", want does not match",
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			t.Run("Error", func(t *testing.T) {
				defer catchFatalPanic(t)
				spy := &testerSpy{}
				check.Error(spy, c.gotErr, c.expectedErr, cmpopts.EquateErrors())
				examineSpyResults(t, spy, c)
			})

			t.Run("ErrorString", func(t *testing.T) {
				defer catchFatalPanic(t)
				spy := &testerSpy{}
				check.ErrorString(spy, c.gotErr, c.expectedMsg)
				examineSpyResults(t, spy, c)
				if spy.fatal != c.wantErrMsgString {
					t.Errorf("ERROR: got %v want %v", spy.fatal, c.wantErrMsgString)
				}
			})
		})
	}
}

// FatalPanic is used to signal a panic that is caused by calling
// Fatalf. It is used to mimic the behavior of testing.T:
//
// Fatalf leads to skipping the remaining code of the test function.
// By executing a panic and catching it in the test call we can achieve
// the same behavior.
//
// Using runtime.Goexit() would lead to an unrecoverable panic of the go test call.
type FatalPanic struct{}

// catchFatalPanic catches the panic that is raised in Fatalf.
func catchFatalPanic(t *testing.T) {
	if p := recover(); p != nil {
		if _, ok := p.(FatalPanic); !ok {
			panic(p) // panic not from Fatalf, rethrow
		}
		t.Logf("INFO: Fatalf skips rest of code")
	}
}

func examineSpyResults(t *testing.T, spy *testerSpy, c struct {
	name             string
	gotErr           error
	expectedMsg      string
	expectedErr      error
	wantErrMsgString string
	wantErrMsg       string
	wantSkipMsg      string
	wantLogMsg       string
}) {

	// Helper() should always be called.
	if !spy.helper {
		t.Errorf("ERROR: helper not called")
	}

	if spy.skip != c.wantSkipMsg {
		t.Errorf("ERROR: got %v, want %v", spy.skip, c.wantSkipMsg)
	}

	if testing.Verbose() {
		if spy.log != c.wantLogMsg {
			t.Errorf("ERROR: got %v, want %v", spy.log, c.wantLogMsg)
		}
	}
}

// testerSpy implements the errStringTester interface.
type testerSpy struct {
	fatal, log, skip string
	helper           bool
}

// Helper tracks if it is called.
func (t *testerSpy) Helper() {
	t.helper = true
}

// Skip tracks the message.
func (t *testerSpy) Skip(info ...any) {
	t.skip = fmt.Sprintf("%v", info)
}

// Fatalf tracks the message and the arguments.
// It panics to simulate functionality of testing.T:
// avoid execution of the steps that follow the call to Fatalf.
func (t *testerSpy) Fatalf(msg string, args ...any) {
	t.fatal = fmt.Sprintf(msg, args...)
	panic(FatalPanic{})
}

// Logf tracks the message and the arguments.
func (t *testerSpy) Logf(msg string, args ...any) {
	t.log = fmt.Sprintf(msg, args...)
}
