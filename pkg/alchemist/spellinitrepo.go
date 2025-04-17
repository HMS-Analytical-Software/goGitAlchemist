package alchemist

import (
	"path/filepath"
	"strings"
)

// initRepoSpell provides initializing a bare repo and cloning it.
type initRepoSpell struct {
	Bare    string `yaml:"bare"`
	CloneTo string `yaml:"clone_to"`
}

// validate checks the values and reports an error if something is missing.
func (s initRepoSpell) validate() error {
	if s.Bare == "" {
		return MissingValueError("bare")
	}
	if s.CloneTo == "" {
		return MissingValueError("clone_to")
	}
	return nil
}

// cast creates a bare repo and clones and initializes it.
func (s initRepoSpell) cast(a assistant, opt Options) error {

	bareDir := filepath.Join(opt.RepoDir, s.Bare)

	a.info("%d/%d: make directory %s", opt.currentSpell, opt.numberOfSpells, bareDir)

	err := a.makedir(bareDir)
	if err != nil {
		return err
	}

	hints := []spellHint{{
		dir:  bareDir,
		args: []string{"init", "--bare", "--initial-branch=" + defaultBranch, "."},
	}, {
		dir:  opt.RepoDir,
		args: []string{"clone", s.Bare, s.CloneTo},
	}, {
		dir:  filepath.Join(opt.RepoDir, s.CloneTo),
		args: []string{"remote", "set-url", "origin", filepath.Join("..", s.Bare)},
	}, {
		dir:  filepath.Join(opt.RepoDir, s.CloneTo),
		args: []string{"config", "user.name", author["red"]},
	}, {
		dir:  filepath.Join(opt.RepoDir, s.CloneTo),
		args: []string{"config", "user.email", email["red"]},
	}, {
		dir:  filepath.Join(opt.RepoDir, s.CloneTo),
		args: []string{"config", "init.defaultBranch", defaultBranch},
	}}

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
