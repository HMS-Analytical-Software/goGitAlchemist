package alchemist

import "path/filepath"

// pushSpell provides push to a remote repository.
type pushSpell struct {
	Main bool `yaml:"main"`
}

// validate checks the values and reports an error if something is missing.
func (s pushSpell) validate() error {
	return nil
}

// cast executes a git push if Main is true.
func (s pushSpell) cast(a assistant, opt Options) error {

	dir := filepath.Join(opt.RepoDir, opt.cloneTo)
	a.info("%d/%d: push (%v)", opt.currentSpell, opt.numberOfSpells, s.Main)

	if s.Main {
		args := []string{"push", "origin", "main"}
		err := a.git(dir, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
