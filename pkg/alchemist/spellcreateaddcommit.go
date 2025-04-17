package alchemist

import (
	"regexp"
)

// createAddCommitSpell
type createAddCommitSpell struct {
	Files   []string `yaml:"files"`
	Message string   `yaml:"message"`
	Author  string   `yaml:"author"`
}

// regexpSplitCreateAddCommit defines the regular  expression for
// splitting the source and target speciffication in the files list.
var regexpSplitCreateAddCommit = regexp.MustCompile(`\s*=>\s*`)

// validate checks the values and reports an error if something is missing.
func (s createAddCommitSpell) validate() error {
	if len(s.Files) == 0 {
		return MissingValueError("files")
	}
	if s.Message == "" {
		return MissingValueError("message")
	}
	if s.Author == "" {
		return MissingValueError("author")
	}

	for _, file := range s.Files {
		elements := regexpSplitCreateAddCommit.Split(file, -1)
		if len(elements) < 2 {
			return InvalidValueError{Variable: "files", Reason: "missing '=>'"}
		}
		if elements[0] == "" {
			return InvalidValueError{Variable: "files", Reason: "source missing"}
		}
		if elements[1] == "" {
			return InvalidValueError{Variable: "files", Reason: "target missing"}
		}
	}

	return nil
}

// cast copies the files, adds them to the index and commits it.
// It uses createFileSpell, addSpell and commitSpell.
func (s createAddCommitSpell) cast(a assistant, opt Options) error {

	spells := make([]caster, 0, len(s.Files)+2)

	for _, filePair := range s.Files {
		elements := regexpSplitCreateAddCommit.Split(filePair, -1)
		spells = append(spells,
			createFileSpell{Source: elements[0], Target: elements[1]})
	}
	spells = append(spells, addSpell{Files: []string{"."}})
	spells = append(spells, commitSpell{
		Message: s.Message,
		Author:  s.Author,
	})

	for _, spell := range spells {
		err := spell.cast(a, opt)
		if err != nil {
			return err
		}
	}

	return nil
}
