package alchemist

import (
	"path/filepath"
	"strings"
)

// gitSpell provides arbitrary git commands.
type gitSpell struct {
	Command string `yaml:"command"`
}

// validate checks the values and reports an error if something is missing.
func (s gitSpell) validate() error {
	if s.Command == "" {
		return MissingValueError("command")
	}
	return nil
}

// cast executes an arbitrary git command.
func (s gitSpell) cast(a assistant, opt Options) error {

	dir := filepath.Join(opt.RepoDir, opt.cloneTo)
	args := s.splitArgs()
	a.info("%d/%d: %s", opt.currentSpell, opt.numberOfSpells, strings.Join(args, " "))

	err := a.git(dir, args...)
	if err != nil {
		return err
	}

	return nil
}

// splitArgs splits args like the shell: double quoted text
// can contain spaces.
func (s gitSpell) splitArgs() []string {

	var quoted bool
	args := strings.FieldsFunc(s.Command, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
			return true
		}
		return !quoted && r == ' '
	})

	// strip prefixed git command
	if len(args) > 0 && (args[0] == linuxGitCmd || args[0] == windowsGitCmd) {
		args = args[1:]
	}

	return args
}
