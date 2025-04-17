package alchemist

import (
	"errors"
	"testing"
)

// TestErrorMessage tests the error types.
func TestErrorMessage(t *testing.T) {

	elementErrMsg := "invalid element"
	execErrMsg := "exit status 42"

	elementErr := errors.New(elementErrMsg)
	execErr := errors.New(execErrMsg)

	testCases := []struct {
		name string
		err  error
		want string
	}{{
		name: "MissingValueError",
		err:  MissingValueError("x"),
		want: "value for x is missing",
	}, {
		name: "YamlDecodeError",
		err:  YamlDecodeError{Element: "node", Err: elementErr},
		want: "yaml decode node: " + elementErrMsg,
	}, {
		name: "unwrap YamlDecodeError",
		err:  YamlDecodeError{Err: elementErr}.Unwrap(),
		want: elementErrMsg,
	}, {
		name: "InvalidValueError",
		err:  InvalidValueError{Variable: "x", Reason: "y"},
		want: "value for x: y",
	}, {
		name: "ExecError",
		err:  ExecError{Cmd: "git", Args: []string{"x", "y", "z"}, Err: execErr},
		want: "git x y z: " + execErrMsg,
	}, {
		name: "unwrap ExecError",
		err: ExecError{
			Cmd:  "git",
			Args: []string{"x", "y", "z"},
			Err:  execErr,
		}.Unwrap(),
		want: execErrMsg,
	}, {
		name: "IOError",
		err:  IOError{Cmd: "open", Arg: "file", Err: execErr},
		want: "open: " + execErrMsg,
	}, {
		name: "Unwrap IOError",
		err:  IOError{Cmd: "open", Arg: "file", Err: execErr}.Unwrap(),
		want: execErrMsg,
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			got := c.err.Error()
			if got != c.want {
				t.Errorf("ERROR: got: %v, want: %v", got, c.want)
			}
		})
	}
}
