//go:build !acctest

package main

import (
	"errors"
	"io"
	"log"
	"path/filepath"
	"testing"

	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist"
	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/check"
)

const testDataDir = "testdata"

var transmuteSpyError = errors.New("transmuteSpy error")

type transmuteSpy struct {
	formula    alchemist.Formula
	opt        alchemist.Options
	raiseError bool
}

func (s *transmuteSpy) transmute(formula alchemist.Formula, opt alchemist.Options, _ *log.Logger) error {
	s.formula = formula
	s.opt = opt
	if s.raiseError {
		return transmuteSpyError
	}
	return nil
}

func TestRunOneTask(t *testing.T) {

	testCases := []struct {
		name                  string
		fileName              string
		opt                   alchemist.Options
		spy                   *transmuteSpy
		wantPath, wantFormula string
		wantErr               string
	}{{
		name: "ok",
		fileName: filepath.Join(testDataDir, "basic_workflow",
			alchemist.FormulaFileName),
		opt:         alchemist.Options{CfgDir: testDataDir},
		spy:         &transmuteSpy{},
		wantPath:    "basic_workflow",
		wantFormula: "basic_workflow",
	}, {
		name: "no config dir",
		fileName: filepath.Join(testDataDir, "basic_workflow",
			alchemist.FormulaFileName),
		spy:         &transmuteSpy{},
		wantPath:    filepath.Join(testDataDir, "basic_workflow"),
		wantFormula: "basic_workflow",
	}, {
		name: "formula title not matching path",
		fileName: filepath.Join(testDataDir, "cmd_merge",
			alchemist.FormulaFileName),
		opt:         alchemist.Options{CfgDir: testDataDir},
		spy:         &transmuteSpy{},
		wantPath:    "cmd_merge",
		wantFormula: "merge",
	}, {
		name:     "file not found",
		fileName: "file not found",
		spy:      &transmuteSpy{},
		wantErr:  "open file not found: ",
	}, {
		name: "transmute error",
		fileName: filepath.Join(testDataDir, "basic_workflow",
			alchemist.FormulaFileName),
		opt:     alchemist.Options{CfgDir: testDataDir},
		spy:     &transmuteSpy{raiseError: true},
		wantErr: "transmuteSpy error",
	}, {
		name: "rel error",
		fileName: filepath.Join(testDataDir, "basic_workflow",
			alchemist.FormulaFileName),
		opt: alchemist.Options{CfgDir: "/"},
		spy: &transmuteSpy{},
		wantErr: "can't make " + filepath.Join(testDataDir, "basic_workflow") +
			" relative to /",
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := runOneTask(c.spy.transmute, c.fileName, c.opt,
				log.New(io.Discard, "", log.LstdFlags))
			check.ErrorString(t, err, c.wantErr)
			if c.spy.formula.Title != c.wantFormula {
				t.Errorf("ERROR: got: %v, want: %v", c.spy.formula.Title, c.wantFormula)
			}
			if c.spy.opt.TaskDir != c.wantPath {
				t.Errorf("ERROR: got: %v, want: %v", c.spy.opt.TaskDir, c.wantPath)

			}
		})
	}

}

func TestRunTaskList(t *testing.T) {

	testCases := []struct {
		name     string
		fileList []string
		spy      *transmuteSpy
		wantErr  string
	}{{
		name: "ok",
		fileList: []string{
			filepath.Join(testDataDir, "basic_workflow", alchemist.FormulaFileName),
			filepath.Join(testDataDir, "cmd_merge", alchemist.FormulaFileName),
		},
		spy: &transmuteSpy{},
	}, {
		name: "error",
		fileList: []string{
			filepath.Join(testDataDir, "basic_workflow", alchemist.FormulaFileName),
			filepath.Join(testDataDir, "cmd_merge", alchemist.FormulaFileName),
		},
		spy:     &transmuteSpy{raiseError: true},
		wantErr: "transmuteSpy error",
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			err := runTaskList(c.spy.transmute, c.fileList, alchemist.Options{},
				log.New(io.Discard, "", log.LstdFlags))
			check.ErrorString(t, err, c.wantErr)
		})
	}
}
