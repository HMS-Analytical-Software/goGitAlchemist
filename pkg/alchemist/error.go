package alchemist

import "strings"

// MissingValueError signals a missing value in the recipe definition.
type MissingValueError string

// Error implents the error interface.
func (e MissingValueError) Error() string {
	return "value for " + string(e) + " is missing"
}

// InvalidValueError signals an invalid value in the recipe definition.
type InvalidValueError struct {
	Variable, Reason string
}

// Error implents the error interface.
func (e InvalidValueError) Error() string {
	return "value for " + e.Variable + ": " + e.Reason
}

// YamlDecodeError signals an error during yaml decoding.
// It keeps the underlying error that was returned by the yaml package.
type YamlDecodeError struct {
	Element string
	Err     error
}

// Error implents the error interface.
func (e YamlDecodeError) Error() string {
	return "yaml decode " + e.Element + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e YamlDecodeError) Unwrap() error {
	return e.Err
}

// ExecError signals an error that happend during execution.
// It keeps the underlying error that was returned by the os package.
type ExecError struct {
	Cmd  string
	Args []string
	Err  error
}

// Error returns the command, the arguments and the message
// of the underlying error.
// It implents the error interface.
func (e ExecError) Error() string {
	var result string
	if e.Cmd != "" {
		result = e.Cmd + " "
	}
	if len(e.Args) > 0 {
		result += strings.Join(e.Args, " ")
	}
	if e.Cmd != "" {
		result += ": "
	}

	return result + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e ExecError) Unwrap() error {
	return e.Err
}

// IOError signals an error that happend during i/o operations.
// It keeps the underlying error that was returned by the os package.
type IOError struct {
	Cmd string
	Arg string
	Err error
}

// Error returns the command and the message from the os which usually contains
// all relevant info.
func (e IOError) Error() string {
	return e.Cmd + ": " + e.Err.Error()
}

// Unwrap returns the underlying error.
func (e IOError) Unwrap() error {
	return e.Err
}
