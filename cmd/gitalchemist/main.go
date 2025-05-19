// gitalchemist provides executing gitalchemist formulas.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist"
)

func main() {

	opt, err := getOptions(os.Args, os.Getenv, os.Stderr)
	if err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		os.Exit(1)
	}

	if opt.version {
		fmt.Println(Version)
		return
	}

	err = run(opt)
	if err != nil {
		log.Printf("[ERROR] %v", err)

		// return code depending on error type
		var missingErr alchemist.MissingValueError
		if errors.As(err, &missingErr) {
			os.Exit(2)
		}
		var invalidErr alchemist.InvalidValueError
		if errors.As(err, &invalidErr) {
			os.Exit(3)
		}
		var yamlErr alchemist.YamlDecodeError
		if errors.As(err, &yamlErr) {
			os.Exit(4)
		}
		var execErr alchemist.ExecError
		if errors.As(err, &execErr) {
			os.Exit(5)
		}
		var ioErr alchemist.IOError
		if errors.As(err, &ioErr) {
			os.Exit(6)
		}

		os.Exit(42)
	}

	log.Printf("[INFO] ok")
}

// options represent the settings from the command line.
type options struct {
	targetdir string
	cfgDir    string
	verbose   bool
	test      bool
	taskList  []string
	maxSteps  int
	runAll    bool
	clean     bool
	version   bool
}

// run executes the main program and returns the error status.
func run(opt options) (err error) {

	// clean output directory
	if opt.clean {
		log.Printf("[INFO] remove %s", opt.targetdir)
		return os.RemoveAll(opt.targetdir)
	}

	alchemistOpt := alchemist.Options{
		RepoDir:       opt.targetdir,
		CfgDir:        opt.cfgDir,
		Verbose:       opt.verbose,
		Test:          opt.test,
		ExecuteSpells: opt.maxSteps,
	}

	var fileList []string
	if opt.runAll {
		fileList, err = alchemist.ListBookContent(alchemistOpt.CfgDir)
	} else {
		fileList, err = alchemist.ListPages(alchemistOpt.CfgDir, opt.taskList...)
	}
	if err != nil {
		return err
	}

	logger := log.New(os.Stderr, "", log.LstdFlags)
	return runTaskList(alchemist.Transmute, fileList, alchemistOpt, logger)
}

// transmuteFunc defines the function type for transmuting matter.
type transmuteFunc func(alchemist.Formula, alchemist.Options, *log.Logger) error

// runTaskList runs all the tasks in the list.
// If it gets an error, it returns it immediately.
func runTaskList(fn transmuteFunc, fileList []string, opt alchemist.Options, logger *log.Logger) error {
	for _, file := range fileList {
		err := runOneTask(fn, file, opt, logger)
		if err != nil {
			return err
		}
	}
	return nil
}

// runOneTask executes one gitalchemist formula.
func runOneTask(fn transmuteFunc, file string, opt alchemist.Options, logger *log.Logger) error {
	// runOneTask executes one gitalchemist formula.
	formula, err := alchemist.Read(file)
	if err != nil {
		return err
	}
	path, err := filepath.Rel(opt.CfgDir, filepath.Dir(file))
	if err != nil {
		return err
	}
	opt.TaskDir = path

	return fn(formula, opt, logger)
}

// getOptions retrieves the command line options.
// Providing the information as parameters facilitates testing.
// - args: the command line arguments
// - getenv: the environment variables
// - stderr: the target for messages
func getOptions(args []string, getenv getEnvFunc, stderr io.Writer) (options, error) {

	var f flag.FlagSet
	f.SetOutput(stderr)

	// define help text
	f.Usage = func() {
		fmt.Fprintf(stderr, `usage: %s <path/to/dir> [ path/to/dir ... ]]
usage: %[1]s -runall
usage: %[1]s -clean

The directories must contain a definition file named %q 
and all the files that are used in the definition.

`, args[0], alchemist.FormulaFileName)
		fmt.Fprintf(stderr, "usage of %s:\n", args[0])
		f.PrintDefaults()
	}

	// define flags
	targetDir := getenv("GITALCHEMIST_TARGETDIR")
	if targetDir == "" {
		targetDir = defaultCwd
	}
	optTargetDir := f.String("targetdir", targetDir,
		"base directory for generated git repos (default: $GITALCHEMIST_TARGETDIR)")
	optCfgDir := f.String("cfgdir", getenv("GITALCHEMIST_CFGDIR"),
		"base directory for git alchemy recipes (default: $GITALCHEMIST_CFGDIR)")
	optSteps := f.Int("maxsteps", 0, "execute only number of specified steps\n"+
		"0 executes all steps")

	optVerbose := f.Bool("verbose", false, "verbose messages")
	optTest := f.Bool("test", false, "test run, steps are logged but not executed")
	optRunAll := f.Bool("runall", false, "run all recipes")
	optClean := f.Bool("clean", false, "remove targetdir")
	optVersion := f.Bool("version", false, "show version")

	// parse command line flags
	err := f.Parse(args[1:])
	if err != nil {
		return options{}, err
	}

	taskList := f.Args()

	// check mandatory arguments
	count := 0
	if *optClean {
		count++
	}
	if *optRunAll {
		count++
	}
	if len(taskList) > 0 {
		count++
	}
	if count != 1 && !*optVersion {
		f.Usage()
		return options{}, fmt.Errorf("specify a task, runall, or clean")
	}

	return options{
		targetdir: *optTargetDir,
		cfgDir:    *optCfgDir,
		verbose:   *optVerbose,
		test:      *optTest,
		taskList:  taskList,
		maxSteps:  *optSteps,
		runAll:    *optRunAll,
		clean:     *optClean,
		version:   *optVersion,
	}, nil
}

// Version contains the git tag this binary was built with.
// It is set during compilation on the command line.
var Version string

// getEnvFunc is used to get environment variable values.
// It is used instead of os.Getenv to facilitate unit testing.
type getEnvFunc func(string) string

const (
	// defaultCwd defines the default base dir for all git repos.
	defaultCwd = "cwd"
)
