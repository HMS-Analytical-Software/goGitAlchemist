package alchemist

import (
	"path/filepath"
	"strings"
)

// mergeSpell provides merging two git branches.
type mergeSpell struct {
	Source       string `yaml:"source"`
	Target       string `yaml:"target"`
	DeleteSource bool   `yaml:"delete_source"`
}

// validate checks the values and reports an error if something is missing.
func (s mergeSpell) validate() error {
	if s.Source == "" {
		return MissingValueError("source")
	}
	if s.Target == "" {
		return MissingValueError("target")
	}
	return nil
}

// cast executes a git merge.
func (s mergeSpell) cast(a assistant, opt Options) error {

	dir := filepath.Join(opt.RepoDir, opt.cloneTo)
	a.info("%d/%d: merge %s with %s  (delete: %v)",
		opt.currentSpell, opt.numberOfSpells, s.Source, s.Target, s.DeleteSource)

	hints := []spellHint{{
		dir:  dir,
		args: []string{"checkout", s.Target},
	}, {
		dir:  dir,
		args: []string{"merge", s.Source},
	}}

	if s.DeleteSource {
		hints = append(hints, spellHint{
			dir:  dir,
			args: []string{"branch", "-d", s.Source},
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

	return nil
}
