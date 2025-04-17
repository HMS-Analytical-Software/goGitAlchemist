package alchemist

import (
	"fmt"
	"strings"
)

// assistantSpy allows recording calls in unit tests.
// It implements the assistant interface.
type assistantSpy struct {
	counter int        // track number of calls
	errorAt int        // return error at this call, count starts with 1
	calls   [][]string // recorded calls

	mortalLogger // noop, just to implement assistant interface
}

// git tracks the git calls.
// If errorAt is reached, an error is returned.
func (s *assistantSpy) git(dir string, args ...string) error {
	s.counter++
	if s.counter == s.errorAt {
		return fmt.Errorf("%s %s: spy error: %d", gitCmd, strings.Join(args, " "),
			s.counter)
	}
	s.calls = append(s.calls, append([]string{dir, gitCmd}, args...))
	return nil
}

// copy tracks the copy calls.
// If errorAt is reached, an error is returned.
func (s *assistantSpy) copy(from, to string) error {
	s.counter++
	if s.counter == s.errorAt {
		return fmt.Errorf("copy %s %s: spy error: %d", from, to, s.counter)
	}
	s.calls = append(s.calls, []string{"copy", from, to})
	return nil
}

// makedir tracks the makedir calls.
// If errorAt is reached, an error is returned.
func (s *assistantSpy) makedir(dir string) error {
	s.counter++
	if s.counter == s.errorAt {
		return fmt.Errorf("make dir %s: spy error: %d", dir, s.counter)
	}
	s.calls = append(s.calls, append([]string{"makedir"}, dir))
	return nil
}
