package alchemist

import (
	"errors"
	"testing"
)

// TestCasterValidate tests the validate function of all casters.
func TestCasterValidate(t *testing.T) {

	testCases := []struct {
		name    string
		spell   caster
		wantErr error
	}{{
		name:  "initRepoSpell ok",
		spell: initRepoSpell{Bare: "x", CloneTo: "y"},
	}, {
		name:    "initRepoSpell bare missing",
		spell:   initRepoSpell{CloneTo: "y"},
		wantErr: MissingValueError("bare"),
	}, {
		name:    "initRepoSpell clone_to missing",
		spell:   initRepoSpell{Bare: "x"},
		wantErr: MissingValueError("clone_to"),
	}, {
		name:    "initRepoSpell both missing",
		spell:   initRepoSpell{},
		wantErr: MissingValueError("bare"),
	}, {
		name:  "createFileSpell ok",
		spell: createFileSpell{Source: "x", Target: "y"},
	}, {
		name:    "createFileSpell source missing",
		spell:   createFileSpell{Target: "y"},
		wantErr: MissingValueError("source"),
	}, {
		name:    "createFileSpell target missing",
		spell:   createFileSpell{Source: "x"},
		wantErr: MissingValueError("target"),
	}, {
		name:  "addSpell ok",
		spell: addSpell{Files: []string{"x"}},
	}, {
		name:    "addSpell files missing",
		spell:   addSpell{Files: []string{}},
		wantErr: MissingValueError("files"),
	}, {
		name:    "addSpell files nil",
		spell:   addSpell{},
		wantErr: MissingValueError("files"),
	}, {
		name:  "commitSpell ok",
		spell: commitSpell{Message: "x", Author: "y"},
	}, {
		name:    "commitSpell message missing",
		spell:   commitSpell{Author: "y"},
		wantErr: MissingValueError("message"),
	}, {
		name:    "commitSpell author missing",
		spell:   commitSpell{Message: "x"},
		wantErr: MissingValueError("author"),
	}, {
		name:  "createAddCommitSpell ok",
		spell: createAddCommitSpell{Files: []string{"x=>a"}, Message: "y", Author: "z"},
	}, {
		name:    "createAddCommitSpell separator missing",
		spell:   createAddCommitSpell{Files: []string{"x-a"}, Message: "y", Author: "z"},
		wantErr: InvalidValueError{Variable: "files", Reason: "missing '=>'"},
	}, {
		name:    "createAddCommitSpell source missing",
		spell:   createAddCommitSpell{Files: []string{"=>a"}, Message: "y", Author: "z"},
		wantErr: InvalidValueError{Variable: "files", Reason: "source missing"},
	}, {
		name:    "createAddCommitSpell target missing",
		spell:   createAddCommitSpell{Files: []string{"x=>"}, Message: "y", Author: "z"},
		wantErr: InvalidValueError{Variable: "files", Reason: "target missing"},
	}, {
		name:    "createAddCommitSpell files missing",
		spell:   createAddCommitSpell{Message: "y", Author: "z"},
		wantErr: MissingValueError("files"),
	}, {
		name:    "createAddCommitSpell message missing",
		spell:   createAddCommitSpell{Files: []string{"x"}, Author: "z"},
		wantErr: MissingValueError("message"),
	}, {
		name:    "createAddCommitSpell author missing",
		spell:   createAddCommitSpell{Files: []string{"x"}, Message: "y"},
		wantErr: MissingValueError("author"),
	}, {
		name:    "createAddCommitSpell all missing",
		spell:   createAddCommitSpell{},
		wantErr: MissingValueError("files"),
	}, {
		name:  "gitSpell ok",
		spell: gitSpell{Command: "x"},
	}, {
		name:    "gitSpell command missing",
		spell:   gitSpell{},
		wantErr: MissingValueError("command"),
	}, {
		name:  "mergeSpell ok",
		spell: mergeSpell{Source: "x", Target: "y"},
	}, {
		name:    "mergeSpell source missing",
		spell:   mergeSpell{Target: "y"},
		wantErr: MissingValueError("source"),
	}, {
		name:    "mergeSpell target missing",
		spell:   mergeSpell{Source: "x"},
		wantErr: MissingValueError("target"),
	}, {
		name:  "pushSpell ok",
		spell: pushSpell{},
	}, {
		name:  "moveSpell ok",
		spell: moveSpell{Source: "x", Target: "y"},
	}, {
		name:    "moveSpell source missing",
		spell:   moveSpell{Target: "y"},
		wantErr: MissingValueError("source"),
	}, {
		name:    "moveSpell target missing",
		spell:   moveSpell{Source: "x"},
		wantErr: MissingValueError("target"),
	}, {
		name:  "removeAndCommitSpell ok",
		spell: removeAndCommitSpell{Files: []string{"x"}, Message: "y", Author: "z"},
	}, {
		name:    "removeAndCommitSpell files missing",
		spell:   removeAndCommitSpell{Message: "y", Author: "z"},
		wantErr: MissingValueError("files"),
	}, {
		name:    "removeAndCommitSpell message missing",
		spell:   removeAndCommitSpell{Files: []string{"x"}, Author: "z"},
		wantErr: MissingValueError("message"),
	}, {
		name:    "removeAndCommitSpell author missing",
		spell:   removeAndCommitSpell{Files: []string{"x"}, Message: "y"},
		wantErr: MissingValueError("author"),
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			err := c.spell.validate()

			if c.wantErr != nil {
				if !errors.Is(err, c.wantErr) {
					t.Errorf("ERROR: got: %v, want: %v", err, c.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ERROR: got error: %v", err)
			}
		})
	}
}
