package alchemist

import (
	"path/filepath"
	"strings"
)

// removeAndCommitSpell provides the combined command of removing files
// and commiting them.
type removeAndCommitSpell struct {
	Files   []string `yaml:"files"`
	Message string   `yaml:"message"`
	Author  string   `yaml:"author"`
}

// validate checks the values and reports an error if something is missing.
func (s removeAndCommitSpell) validate() error {
	if len(s.Files) == 0 {
		return MissingValueError("files")
	}
	if s.Message == "" {
		return MissingValueError("message")
	}
	if s.Author == "" {
		return MissingValueError("author")
	}

	return nil
}

// cast calls git rm to all the files and commits the result.
func (s removeAndCommitSpell) cast(a assistant, opt Options) error {

	dir := filepath.Join(opt.RepoDir, opt.cloneTo)
	a.info("%d/%d: remove and commit %d files", opt.currentSpell, opt.numberOfSpells,
		len(s.Files))

	hints := make([]spellHint, 0, len(s.Files))

	for _, fileName := range s.Files {
		hints = append(hints, spellHint{
			dir:  dir,
			args: []string{"rm", fileName},
		})
	}
	for _, hint := range hints {
		a.info("%d/%d: %s", opt.currentSpell, opt.numberOfSpells,
			strings.Join(hint.args, " "))
		err := a.git(hint.dir, hint.args...)
		if err != nil {
			return err
		}
	}

	return commitSpell{Message: s.Message, Author: s.Author}.cast(a, opt)
}
