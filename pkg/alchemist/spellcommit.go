package alchemist

import "path/filepath"

// commitSpell provides committing the index to the repo.
type commitSpell struct {
	Message string `yaml:"message"`
	Author  string `yaml:"author"`
}

// validate checks the values and reports an error if something is missing.
func (s commitSpell) validate() error {
	if s.Message == "" {
		return MissingValueError("message")
	}
	if s.Author == "" {
		return MissingValueError("author")
	}
	return nil
}

const gitCommitDateFormat = "format:relative:5.hours.ago"

// incant executes git commit.
func (s commitSpell) cast(a assistant, opt Options) error {

	a.info("%d/%d: commit", opt.currentSpell, opt.numberOfSpells)

	dir := filepath.Join(opt.RepoDir, opt.cloneTo)
	args := []string{"commit", "--date=" + gitCommitDateFormat,
		"-m", s.Message, "--author=" + getAuthor(s.Author)}

	err := a.git(dir, args...)
	if err != nil {
		return err
	}

	return nil
}
