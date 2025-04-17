package alchemist

import (
	"log"
)

// novice is an assistant that does not execute the commands, it just
// reports the steps that should be executed (on debug level).
//
// It implements the assistant interface.
//
// It is used for test mode. All methods return a nil error.
type novice struct {
	mortalLogger
}

// newNovice returns an initialized novice object.
func newNovice(l *log.Logger, opt Options) novice {
	return novice{
		mortalLogger: mortalLogger{
			Logger:  l,
			verbose: opt.Verbose,
		},
	}
}

// git emits a debug message with the parameters.
// It implements the assistant interface.
func (n novice) git(dir string, args ...string) error {
	n.debug("%q: git %#v", dir, args)
	return nil
}

// copy emits a debug message with the parameters.
// It implements the assistant interface.
func (n novice) copy(from, to string) error {
	n.debug("copy %q to %q", from, to)
	return nil
}

// makedir emits a debug message with the parameters.
// It implements the assistant interface.
func (n novice) makedir(dir string) error {
	n.debug("makedir %q", dir)
	return nil
}
