package alchemist

import "path/filepath"

// addSpell  provides adding files to the git index.
type addSpell struct {
	Files []string `yaml:"files"`
}

// validate checks the values and reports an error if something is missing.
func (s addSpell) validate() error {
	if len(s.Files) == 0 {
		return MissingValueError("files")
	}
	return nil
}

// cast executes git add of all files in the list.
func (s addSpell) cast(a assistant, opt Options) error {

	a.info("%d/%d: add %d files", opt.currentSpell, opt.numberOfSpells, len(s.Files))

	dir := filepath.Join(opt.RepoDir, opt.cloneTo)
	for _, fileName := range s.Files {

		args := []string{"add", fileName}
		err := a.git(dir, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
