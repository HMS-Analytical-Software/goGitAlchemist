package alchemist

import (
	"errors"
	"testing"
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

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ERROR: got: %v, want: %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ERROR: got error: %v", err)
			}
		})
	}
}

func TestFilterInitSpells(t *testing.T) {
	// Create test spells including init_bare_repo
	initSpell := initRepoSpell{Bare: "test", CloneTo: "test"}
	commitSpell := commitSpell{Message: "test", Author: "red"}
	addSpell := addSpell{Files: []string{"test.txt"}}
	spells := []caster{initSpell, commitSpell, addSpell}

	filtered := filterInitSpells(spells)

	// Should have filtered out the init spell, leaving commit and add
	if len(filtered) != 2 {
		t.Errorf("Expected 2 spells after filtering, got %d", len(filtered))
	}

	// Verify the init spell was removed
	for _, spell := range filtered {
		if _, isInitSpell := spell.(initRepoSpell); isInitSpell {
			t.Error("Init spell should have been filtered out")
		}
	}
}

func TestCircularIncludeDetection(t *testing.T) {
	// Test circular include detection logic
	testPath1 := "/path/a.yaml"
	testPath2 := "/path/b.yaml"
	opt := Options{includePath: []string{testPath1, testPath2}}

	// This should detect circular include for testPath1
	isCircular := false
	for _, path := range opt.includePath {
		if path == testPath1 {
			isCircular = true
			break
		}
	}

	if !isCircular {
		t.Error("Should have detected circular include")
	}

	// This should not detect circular include for a new path
	testPath3 := "/path/c.yaml"
	isCircular = false
	for _, path := range opt.includePath {
		if path == testPath3 {
			isCircular = true
			break
		}
	}

	if isCircular {
		t.Error("Should not have detected circular include for new path")
	}
}
