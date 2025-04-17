package alchemist

import (
	"path/filepath"
)

// createFileSpell provides copying a file to the git repo directory.
type createFileSpell struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

// validate checks the values and reports an error if something is missing.
func (s createFileSpell) validate() error {
	if s.Source == "" {
		return MissingValueError("source")
	}
	if s.Target == "" {
		return MissingValueError("target")
	}
	return nil
}

// cast copies the file to the repo.
func (s createFileSpell) cast(a assistant, opt Options) error {

	from := filepath.Join(opt.CfgDir, opt.TaskDir, s.Source)
	to := filepath.Join(opt.RepoDir, opt.cloneTo, s.Target)
	a.info("%d/%d: copy %s to %s", opt.currentSpell, opt.numberOfSpells, from, to)

	err := a.copy(from, to)
	if err != nil {
		return err
	}
	return nil
}
