package alchemist

import (
	"fmt"
	"path/filepath"
)

// includeSpell provides including and executing another formula.
type includeSpell struct {
	Path string `yaml:"path"`
}

// validate checks the values and reports an error if something is missing.
func (s includeSpell) validate() error {
	if s.Path == "" {
		return MissingValueError("path")
	}
	return nil
}

// cast executes the included formula with inherited context.
func (s includeSpell) cast(a assistant, opt Options) error {
	a.info("%d/%d: include %s", opt.currentSpell, opt.numberOfSpells, s.Path)

	// Resolve path relative to current formula's directory
	// TaskDir is relative to CfgDir, so we need to resolve relative to the full path
	formulaDir := filepath.Join(opt.CfgDir, opt.TaskDir)
	includePath := filepath.Join(formulaDir, s.Path)

	// Check for circular includes
	for _, path := range opt.includePath {
		if path == includePath {
			return fmt.Errorf("circular include detected: %s", includePath)
		}
	}

	// Load the included formula
	formula, err := Read(includePath)
	if err != nil {
		return fmt.Errorf("include %s: %w", s.Path, err)
	}

	// Execute included formula with inherited context
	return executeIncludedFormula(formula, includePath, a, opt)
}

// executeIncludedFormula processes an included formula with context inheritance.
func executeIncludedFormula(formula Formula, formulaPath string, a assistant, parentOpt Options) error {
	// Create new options for included formula
	includeOpt := parentOpt

	// formulaPath is the full path, we need to make TaskDir relative to CfgDir
	fullFormulaDir := filepath.Dir(formulaPath)
	relativeTaskDir, err := filepath.Rel(parentOpt.CfgDir, fullFormulaDir)
	if err != nil {
		return fmt.Errorf("failed to resolve relative path for included formula: %w", err)
	}

	includeOpt.TaskDir = relativeTaskDir
	includeOpt.taskName = formula.Title
	includeOpt.includePath = append(parentOpt.includePath, formulaPath)
	// Note: cloneTo is inherited from parent (same git repository)

	a.info("executing included formula: %s", formula.Title)

	// Filter out init_bare_repo spells and process remaining spells
	filteredSpells := filterInitSpells(formula.Commands.spells)
	includeOpt.numberOfSpells = len(filteredSpells)

	// Execute filtered spells
	for i, spell := range filteredSpells {
		includeOpt.currentSpell = i + 1

		err := spell.cast(a, includeOpt)
		if err != nil {
			return fmt.Errorf("included formula %s: %w", formula.Title, err)
		}
	}

	return nil
}

// filterInitSpells removes init_bare_repo spells from included formulas.
func filterInitSpells(spells []caster) []caster {
	var filtered []caster
	for _, spell := range spells {
		// Skip init_bare_repo spells in included formulas
		if _, isInitSpell := spell.(initRepoSpell); !isInitSpell {
			filtered = append(filtered, spell)
		}
	}
	return filtered
}
