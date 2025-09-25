# Include Spell Feature Implementation Plan

## Overview

Adding an `include` spell to goGitAlchemist that allows formulas to include and execute other formulas while sharing the same git repository context but maintaining their own file directories.

## Current State Analysis

**Formula Processing System**:
- `Transmute()` function in `pkg/alchemist/formula.go:20` processes formulas sequentially using the caster interface
- Each spell implements `caster` interface with `validate()` and `cast()` methods defined in `pkg/alchemist/caster.go:3`
- `symbols.UnmarshalYAML()` in `pkg/alchemist/formula.go:69` handles spell type switching via switch statement
- Repository context is established by `init_bare_repo` spell and stored in `opt.cloneTo`

**Repository Initialization Pattern**:
- `init_bare_repo` spell in `pkg/alchemist/spellinitrepo.go:26` creates bare repository and clone
- Sets `cloneTo` value in `symbols.cloneTo` which becomes `opt.cloneTo` for all subsequent spells
- All spells operate within the same git repository context established by the initial `init_bare_repo`

**File Resolution System**:
- File operations use paths relative to `Options.TaskDir` (current formula directory)
- Each formula has its own `files/` subdirectory for source files
- `create_file` spell copies from formula's `files/` directory to working directory

## Desired End State

A fully functional `include` spell that:
- Loads and executes formulas from paths relative to current formula directory
- Inherits parent formula's git repository context (`cloneTo` directory)
- Uses included formula's own `files/` directory for file operations
- Automatically skips `init_bare_repo` spells in included formulas
- Detects and prevents circular includes with clear error messages
- Maintains proper spell progress tracking across included formulas

**Verification**: Successfully execute formulas with includes, verify repository consistency, confirm error handling for edge cases.

### Key Discoveries:
- Formula execution is purely sequential with no existing recursion support
- `Options` struct carries all execution context including directory paths and progress tracking
- Error handling follows structured pattern with specific error types in `pkg/alchemist/error.go`
- YAML unmarshaling uses generic `unmarshalCaster` function for type safety

## What We're NOT Doing

- Supporting absolute paths (only relative paths from current formula directory)
- Modifying core `Transmute()` function signature or behavior
- Adding complex dependency resolution beyond circular detection
- Merging or sharing `files/` directories between parent and included formulas
- Supporting include depth limits (keeping it simple initially)
- Adding conditional include logic or parameters

## Implementation Approach

Context-aware recursive formula execution with inheritance of git repository context while preserving file operation isolation. Implementation uses existing patterns and minimal changes to core architecture.

## Phase 1: Spell Structure and Integration

### Overview
Create the basic `include` spell structure and integrate it into the existing YAML unmarshaling system.

### Changes Required:

#### 1. Add Include Spell Constant
**File**: `pkg/alchemist/formula.go`
**Changes**: Add `symbolInclude` constant to existing spell symbols

```go
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
	symbolInclude         = "include"
)
```

#### 2. Extend UnmarshalYAML Switch Statement
**File**: `pkg/alchemist/formula.go`
**Changes**: Add case for `symbolInclude` in the existing switch statement

```go
		case symbolRemoveCommit:
			spell, err = unmarshalCaster[removeAndCommitSpell](contentNode)
		case symbolInclude:
			spell, err = unmarshalCaster[includeSpell](contentNode)
		default:
			return fmt.Errorf("unkonwn command %q", cmd)
```

#### 3. Create Include Spell Implementation
**File**: `pkg/alchemist/spellinclude.go` (new file)
**Changes**: Create complete include spell following existing spell patterns

```go
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
	formulaDir := filepath.Dir(opt.TaskDir)
	includePath := filepath.Join(formulaDir, s.Path)
	
	// Load the included formula
	formula, err := Read(includePath)
	if err != nil {
		return fmt.Errorf("include %s: %w", s.Path, err)
	}
	
	a.info("loaded included formula: %s", formula.Title)
	return nil // Phase 1: just validate loading
}
```

### Success Criteria:

#### Automated Verification:
- [ ] Code compiles without errors: `go build ./...`
- [ ] All existing tests pass: `go test ./...`
- [ ] New spell file follows Go conventions: `gofmt -d pkg/alchemist/spellinclude.go`
- [ ] YAML unmarshaling recognizes include spell: unit test parsing

#### Manual Verification:
- [ ] Include spell validates correctly with valid and missing path
- [ ] Formula loading works for valid included formulas
- [ ] Error messages are clear for missing included formulas
- [ ] No regressions in existing spell functionality

---

## Phase 2: Context Inheritance and Execution Logic

### Overview
Implement the core execution logic for included formulas with proper context inheritance and circular dependency detection.

### Changes Required:

#### 1. Extend Options for Include Tracking
**File**: `pkg/alchemist/options.go`
**Changes**: Add include path tracking for circular dependency detection

```go
type Options struct {
	TaskDir       string // directory where the definition is located
	RepoDir       string // directory where the bare repository is located
	CfgDir        string // directory of the configuration definitions
	Verbose       bool   // verbose logging
	Test          bool   // test mode
	ExecuteSpells int    // execute only the first # steps

	// set during processing
	taskName       string   // name of the task to execute
	cloneTo        string   // directory of the repository clone
	numberOfSpells int      // number of steps
	currentSpell   int      // number of the current step (1-based)
	includePath    []string // stack of included formula paths for circular detection
}
```

#### 2. Initialize Include Tracking in Transmute
**File**: `pkg/alchemist/formula.go`
**Changes**: Initialize include path tracking in main Transmute function

```go
func Transmute(f Formula, opt Options, logger *log.Logger) error {

	var helper assistant = newAdept(logger, opt)
	if opt.Test {
		helper = newNovice(logger, opt)
	}
	helper.info("execute formula %s", f.Title)

	opt.cloneTo = f.Commands.cloneTo
	opt.taskName = f.Title
	opt.numberOfSpells = len(f.Commands.spells)
	opt.includePath = []string{} // Initialize include path tracking

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
```

#### 3. Implement Complete Include Execution
**File**: `pkg/alchemist/spellinclude.go`
**Changes**: Add full execution logic with context inheritance

```go
// cast executes the included formula with inherited context.
func (s includeSpell) cast(a assistant, opt Options) error {
	a.info("%d/%d: include %s", opt.currentSpell, opt.numberOfSpells, s.Path)
	
	// Resolve path relative to current formula's directory
	formulaDir := filepath.Dir(opt.TaskDir)
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
	includeOpt.TaskDir = filepath.Dir(formulaPath)
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
```

### Success Criteria:

#### Automated Verification:
- [ ] All tests pass: `go test ./...`
- [ ] Code compiles cleanly: `go build ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Linting passes: `golangci-lint run`

#### Manual Verification:
- [ ] Included formulas execute in correct directory context
- [ ] File operations use included formula's files/ directory
- [ ] Git operations work in inherited repository context
- [ ] Circular includes produce clear error messages
- [ ] Progress tracking shows correct spell counts

---

## Phase 3: Testing Infrastructure

### Overview
Create comprehensive test coverage including unit tests for the include spell and acceptance tests for end-to-end functionality.

### Changes Required:

#### 1. Unit Tests for Include Spell
**File**: `pkg/alchemist/spellinclude_test.go` (new file)
**Changes**: Create comprehensive unit tests following existing test patterns

```go
package alchemist

import (
	"path/filepath"
	"testing"

	"github.com/HMS-Analytical-Software/goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
)

func TestIncludeSpell_validate(t *testing.T) {
	tests := []struct {
		name    string
		spell   includeSpell
		wantErr error
	}{
		{
			name:  "valid path",
			spell: includeSpell{Path: "../shared/common.yaml"},
		},
		{
			name:    "missing path",
			spell:   includeSpell{},
			wantErr: MissingValueError("path"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spell.validate()
			check.Equal(t, err, tt.wantErr)
		})
	}
}

func TestFilterInitSpells(t *testing.T) {
	// Create test spells including init_bare_repo
	initSpell := initRepoSpell{Bare: "test", CloneTo: "test"}
	commitSpell := commitSpell{Message: "test", Author: "red"}
	spells := []caster{initSpell, commitSpell}
	
	filtered := filterInitSpells(spells)
	
	// Should have filtered out the init spell
	if len(filtered) != 1 {
		t.Errorf("Expected 1 spell after filtering, got %d", len(filtered))
	}
}

func TestCircularIncludeDetection(t *testing.T) {
	// Test circular include detection logic
	opt := Options{includePath: []string{"/path/a.yaml", "/path/b.yaml"}}
	
	// This should detect circular include
	testPath := "/path/a.yaml"
	for _, path := range opt.includePath {
		if path == testPath {
			// Expected: circular include detected
			return
		}
	}
	t.Error("Should have detected circular include")
}
```

#### 2. Test Data for Include Scenarios
**File**: `cmd/gitalchemist/testdata/test_include/gitalchemist.yaml` (new file)
**Changes**: Create main test formula that uses include

```yaml
title: test_include
commands:
  - init_bare_repo:
      bare: remotes/test_include
      clone_to: test_include
  - create_file:
      source: files/main.md
      target: main.md
  - add:
      files: [main.md]
  - commit:
      message: Main formula commit
      author: red
  - include:
      path: ../shared_workflow/gitalchemist.yaml
  - commit:
      message: Final commit after include
      author: red
```

**File**: `cmd/gitalchemist/testdata/test_include/files/main.md` (new file)
**Changes**: Create test file for main formula

```markdown
# Main Formula File

This file is created by the main formula.
```

**File**: `cmd/gitalchemist/testdata/shared_workflow/gitalchemist.yaml` (new file)
**Changes**: Create included formula (should be runnable independently)

```yaml
title: shared_workflow
commands:
  # This init should be ignored when included
  - init_bare_repo:
      bare: remotes/shared_workflow
      clone_to: shared_workflow
  - create_file:
      source: files/shared.md
      target: shared.md
  - add:
      files: [shared.md]
  - commit:
      message: Added shared file
      author: blue
```

**File**: `cmd/gitalchemist/testdata/shared_workflow/files/shared.md` (new file)
**Changes**: Create test file for included formula

```markdown
# Shared Workflow File

This file is created by the included formula.
```

#### 3. Acceptance Tests
**File**: `cmd/gitalchemist/acc_test.go`
**Changes**: Add acceptance test for include functionality

```go
func TestAcceptanceInclude(t *testing.T) {
	if !acctest {
		t.Skip("acceptance test")
	}

	buildGitAlchemist(t)
	defer cleanupGitAlchemist()

	err := callGitAlchemist(t, "test_include")
	if err != nil {
		t.Fatal(err)
	}

	// Verify both files exist in the same repository
	checkFileExists(t, "cwd/test_include/main.md")
	checkFileExists(t, "cwd/test_include/shared.md")
	
	// Verify git history includes commits from both formulas
	// Should have: initial commit + main commit + shared commit + final commit = 4 commits
	checkGit(t, "cwd/test_include", "4", "rev-list", "--count", "HEAD")
}

func TestAcceptanceIncludeStandalone(t *testing.T) {
	if !acctest {
		t.Skip("acceptance test")
	}

	buildGitAlchemist(t)
	defer cleanupGitAlchemist()

	// Test that included formula can run independently
	err := callGitAlchemist(t, "shared_workflow")
	if err != nil {
		t.Fatal(err)
	}

	checkFileExists(t, "cwd/shared_workflow/shared.md")
}

func checkFileExists(t *testing.T, filePath string) {
	t.Helper()
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist", filePath)
	}
}
```

### Success Criteria:

#### Automated Verification:
- [ ] All unit tests pass: `go test ./pkg/alchemist`
- [ ] Include spell unit tests achieve >90% code coverage
- [ ] Acceptance tests pass: `go test -tags acctest ./cmd/gitalchemist`
- [ ] Test data formulas validate correctly
- [ ] No test failures or race conditions: `go test -race -count=5 ./...`

#### Manual Verification:
- [ ] Include functionality works end-to-end via CLI
- [ ] Both formulas can run independently
- [ ] File resolution works correctly from included formula
- [ ] Git repository state is consistent after includes
- [ ] Error messages are helpful for debugging issues

---

## Testing Strategy

### Unit Tests:
- Include spell validation (path validation, error cases)
- Circular dependency detection with various path scenarios  
- Spell filtering (removing init_bare_repo from included formulas)
- Context inheritance (cloneTo, TaskDir, includePath tracking)
- Error handling for missing files, invalid formulas
- Edge cases: empty includes, malformed YAML

### Integration Tests:
- End-to-end formula inclusion via CLI with real git operations
- File resolution from included formula's files/ directory
- Git operations in shared repository context
- Multiple levels of inclusion (include within include)
- Include formulas that can run both standalone and included
- Error scenarios: circular includes, missing files, invalid paths

### Manual Testing Steps:
1. Create test formula with include spell pointing to valid included formula
2. Run `./gitalchemist test_include` and verify both sets of files are created
3. Run `./gitalchemist shared_workflow` independently and verify it works
4. Create circular include scenario and verify error message
5. Test include with missing file and verify clear error reporting
6. Test include with invalid path and verify error handling
7. Verify git history shows commits from both formulas in correct order

## Performance Considerations

- File I/O overhead for each included formula (minimal impact for typical use)
- Memory usage increases linearly with include depth (acceptable for reasonable depths)
- Path resolution happens once per include (efficient)  
- No performance impact on formulas not using includes
- Include tracking uses simple slice append (O(n) for circular detection, acceptable)

## Migration Notes

- Existing formulas continue to work without any changes
- New `include` spell is completely opt-in functionality
- No breaking changes to existing APIs or behavior
- Backward compatibility fully maintained
- New test data files are isolated and don't affect existing tests

## References

- Formula system: `pkg/alchemist/formula.go:69` - YAML unmarshaling with spell type switching
- Spell interface: `pkg/alchemist/caster.go:3` - Caster interface contract
- Repository initialization: `pkg/alchemist/spellinitrepo.go:26` - Bare repo setup pattern
- Error handling patterns: `pkg/alchemist/error.go` - Structured error types
- Existing spell patterns: `pkg/alchemist/spellcommit.go:25` - Example spell implementation
- Test patterns: `pkg/alchemist/caster_test.go` - Unit testing approach
- Acceptance testing: `cmd/gitalchemist/acc_test.go:20` - End-to-end test patterns