package alchemist

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Formula contains the instructions from the gitalchemy.yaml file.
type Formula struct {
	Title    string  `yaml:"title"`
	Commands symbols `yaml:"commands"`
}

// Transmute creates the git repositories according to the formula
// and the options. Messages are writte to the logger.
func Transmute(f Formula, opt Options, logger *log.Logger) error {

	var helper assistant = newAdept(logger, opt)
	if opt.Test {
		helper = newNovice(logger, opt)
	}
	helper.info("execute formula %s", f.Title)

	opt.cloneTo = f.Commands.cloneTo
	opt.taskName = f.Title
	opt.numberOfSpells = len(f.Commands.spells)

	for i, spell := range f.Commands.spells {
		if opt.ExecuteSpells > 0 && opt.ExecuteSpells == i {
			return nil
		}
		opt.currentSpell = i + 1

		err := spell.cast(helper, opt)
		if err != nil {
			return err
		}
	}

	return nil
}

// symbols contains a list of spells.
type symbols struct {
	cloneTo string
	spells  []caster
}

// A Formula file can contains these symbols for spells.
const (
	symbolInit            = "init_bare_repo"
	symbolCreateFile      = "create_file"
	symbolAdd             = "add"
	symbolCommit          = "commit"
	symbolGit             = "git"
	symbolCreateAddCommit = "create_add_commit"
	symbolMerge           = "merge"
	symbolPush            = "push"
	symbolMove            = "mv"
	symbolRemoveCommit    = "remove_and_commit"
)

// yaml doku
// gopkg.in/yaml.v3
// https://pkg.go.dev/gopkg.in/yaml.v3

// UnmarshalYAML extracts the commands from the yaml definition
// into a list of casters. It also checks for completeness and
// correctness.
func (c *symbols) UnmarshalYAML(value *yaml.Node) (err error) {

	for i, node := range value.Content {
		if len(node.Content) < 2 {
			return YamlDecodeError{Err: fmt.Errorf("yaml node too small: %#v", node)}
		}

		contentNode := node.Content[1]
		var spell caster

		cmd := node.Content[0].Value
		switch cmd {
		case symbolInit:
			var initData initRepoSpell
			initData, err = unmarshalCaster[initRepoSpell](contentNode)
			// cloneTo is needed in following steps.
			c.cloneTo = initData.CloneTo
			spell = initData
		case symbolCreateFile:
			spell, err = unmarshalCaster[createFileSpell](contentNode)
		case symbolAdd:
			spell, err = unmarshalCaster[addSpell](contentNode)
		case symbolCommit:
			spell, err = unmarshalCaster[commitSpell](contentNode)
		case symbolCreateAddCommit:
			spell, err = unmarshalCaster[createAddCommitSpell](contentNode)
		case symbolGit:
			spell, err = unmarshalCaster[gitSpell](contentNode)
		case symbolMerge:
			spell, err = unmarshalCaster[mergeSpell](contentNode)
		case symbolPush:
			spell, err = unmarshalCaster[pushSpell](contentNode)
		case symbolMove:
			spell, err = unmarshalCaster[moveSpell](contentNode)
		case symbolRemoveCommit:
			spell, err = unmarshalCaster[removeAndCommitSpell](contentNode)
		default:
			return fmt.Errorf("unkonwn command %q", cmd)
		}

		// invalid yaml
		if err != nil {
			return err
		}

		// check for missing values or invalid commands
		err = spell.validate()
		if err != nil {
			return fmt.Errorf("validate %s (%d): %w", cmd, i+1, err)
		}

		c.spells = append(c.spells, spell)
	}

	return nil
}

// unmarshalCaster extracts a single caster from the yaml node.
func unmarshalCaster[T caster](node *yaml.Node) (T, error) {
	var data T
	if err := node.Decode(&data); err != nil {
		return data, YamlDecodeError{
			Element: fmt.Sprintf("node %T", data),
			Err:     err,
		}
	}
	return data, nil
}

// Read reads the file with the definitions and returns the
// content as a Formula object. If an error occurs, it is returned.
func Read(filename string) (Formula, error) {

	file, err := os.Open(filename)
	if err != nil {
		return Formula{}, IOError{Cmd: "open", Arg: filename, Err: err}
	}
	defer file.Close()

	return readFormula(file)
}

// readFormula extracts the yaml content from the provided reader.
func readFormula(r io.Reader) (Formula, error) {
	var result Formula
	err := yaml.NewDecoder(r).Decode(&result)
	if err != nil {
		return result, YamlDecodeError{Element: "Formula", Err: err}
	}
	return result, nil
}
