package alchemist

// assistant define the ability to execute low-level commands.
type assistant interface {
	// git execute a git command in the provided directory
	git(dir string, args ...string) error
	// copy copies a file
	copy(from, to string) error
	// makedir creates a directory
	makedir(dir string) error

	// debug emits a debug message
	debug(msg string, args ...any)
	// info emits an info
	info(msg string, args ...any)
}
