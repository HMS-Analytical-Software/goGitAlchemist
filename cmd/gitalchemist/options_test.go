//go:build !acctest

package main

import (
	"bytes"
	"testing"

	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/check"
	"github.com/google/go-cmp/cmp"
)

const pgmName = "gitalchemist"

func TestGetOptions(t *testing.T) {

	fakeEnv := func(values map[string]string) func(name string) string {
		return func(name string) string {
			return values[name]
		}
	}

	testCases := []struct {
		name    string
		args    []string
		setenv  map[string]string
		want    options
		wantMsg string
		wantErr string
	}{{
		name: "all options from command line",
		args: []string{pgmName,
			"-targetdir", "targetdir",
			"-cfgdir", "cfgdir",
			// "-task", "task1,task2",
			"-maxsteps", "5",
			"-verbose",
			"-test",
			"task1",
			"task2",
		},
		want: options{
			targetdir: "targetdir",
			cfgDir:    "cfgdir",
			verbose:   true,
			test:      true,
			taskList:  []string{"task1", "task2"},
			maxSteps:  5,
		},
	}, {
		name: "values from environment variables",
		setenv: map[string]string{
			"GITALCHEMIST_TARGETDIR": "target_from_env",
			"GITALCHEMIST_CFGDIR":    "cfgdir_from_env",
		},
		args: []string{pgmName, "task1"},
		want: options{
			targetdir: "target_from_env",
			cfgDir:    "cfgdir_from_env",
			taskList:  []string{"task1"},
		},
	}, {
		name: "flags override environment variables",
		setenv: map[string]string{
			"GITALCHEMIST_TARGETDIR": "target_from_env",
			"GITALCHEMIST_CFGDIR":    "cfgdir_from_env",
		},
		args: []string{pgmName,
			"-targetdir", "targetdir_flag",
			"-cfgdir", "cfgdir_flag",
			"task1",
		},
		want: options{
			targetdir: "targetdir_flag",
			cfgDir:    "cfgdir_flag",
			taskList:  []string{"task1"},
		},
	}, {
		name: "runall and default cwd",
		args: []string{pgmName, "-runall"},
		want: options{
			targetdir: defaultCwd,
			runAll:    true,
			taskList:  []string{},
		},
	}, {
		name: "clean and default cwd",
		args: []string{pgmName, "-clean"},
		want: options{
			targetdir: defaultCwd,
			clean:     true,
			taskList:  []string{},
		},
	}, {
		name: "version",
		args: []string{pgmName, "-version"},
		want: options{
			targetdir: defaultCwd,
			version:   true,
			taskList:  []string{},
		},
	}, {
		name:    "task and runall",
		args:    []string{pgmName, "-runall", "task11"},
		wantErr: "specify a task, runall, or clean",
		wantMsg: helpMessage,
	}, {
		name:    "task and clean",
		args:    []string{pgmName, "-clean", "task1"},
		wantErr: "specify a task, runall, or clean",
		wantMsg: helpMessage,
	}, {
		name:    "runall and clean",
		args:    []string{pgmName, "-runall", "-clean"},
		wantErr: "specify a task, runall, or clean",
		wantMsg: helpMessage,
	}, {
		name:    "neither task nor runall",
		args:    []string{pgmName},
		wantErr: "specify a task, runall, or clean",
		wantMsg: helpMessage,
	}, {
		name:    "unknown flag",
		args:    []string{pgmName, "-unknown"},
		wantErr: "flag provided but not defined: -unknown",
		wantMsg: "flag provided but not defined: -unknown\n" + helpMessage,
	}, {
		name:    "help flag",
		args:    []string{pgmName, "-h"},
		wantErr: "flag: help requested",
		wantMsg: helpMessage,
	}}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {

			var buf bytes.Buffer
			if c.setenv == nil {
				c.setenv = map[string]string{}
			}

			gotOpt, err := getOptions(c.args, fakeEnv(c.setenv), &buf)

			// check terminal messages
			if diff := cmp.Diff(buf.String(), c.wantMsg); diff != "" {
				t.Errorf("ERROR: got- want+: %s", diff)
			}

			// check returned error
			check.ErrorString(t, err, c.wantErr)

			// check returned options
			if diff := cmp.Diff(gotOpt, c.want,
				cmp.AllowUnexported(options{}),
			); diff != "" {
				t.Errorf("ERROR: got- want+: %s", diff)
			}
		})
	}
}

var helpMessage = `usage: gitalchemist <path/to/dir> [ path/to/dir ... ]]
usage: gitalchemist -runall
usage: gitalchemist -clean

The directories must contain a definition file named "gitalchemist.yaml" 
and all the files that are used in the definition.

usage of gitalchemist:
  -cfgdir string
    	base directory for git alchemey recipes (default: $GITALCHEMIST_CFGDIR)
  -clean
    	remove targetdir
  -maxsteps int
    	execute only number of specified steps
    	0 executes all steps
  -runall
    	run all recipies
  -targetdir string
    	base directory for generatet git repos (default: $GITALCHEMIST_TARGETDIR) (default "cwd")
  -test
    	test run, steps are logged but not executed
  -verbose
    	verbose messages
  -version
    	show version
`
