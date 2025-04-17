package alchemist

import (
	"path/filepath"
)

// moveSpell provides moving a file in the git clone.
type moveSpell struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

// validate checks the values and reports an error if something is missing.
func (s moveSpell) validate() error {
	if s.Source == "" {
		return MissingValueError("source")
	}
	if s.Target == "" {
		return MissingValueError("target")
	}
	return nil
}

// cast moves a file in the git working directory.
func (s moveSpell) cast(a assistant, opt Options) error {

	dir := filepath.Join(opt.RepoDir, opt.cloneTo)
	a.info("%d/%d: mv %s %s", opt.currentSpell, opt.numberOfSpells,
		s.Source, s.Target)

	args := []string{"mv", s.Source, s.Target}
	err := a.git(dir, args...)
	if err != nil {
		return err
	}

	return nil
}
