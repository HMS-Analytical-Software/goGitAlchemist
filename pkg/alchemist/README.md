# Package alchemist

The structure of the code tries to use the 'alchemy' metaphor for
the naming of types, methods and functions.


#  Formula

A formula is an instruction to create a git repo.
It corresponds to the content of the gitalchemist yaml file.

The formula consists of a series of spells that have to be cast in order.

## caster

A caster is an interface that provides:

* validate: check if all ingredients are prepared
* cast: casting the spell

There are many implementations of caster:

* initRepoSpell: inits a bare repo and clones it.
* createFileSpell: copies a file from the definition area to the git clone directory.
* addSpell: adds files to the git index
* commitSpell: commits the index
* createAddCommitSpell: combines create, add, and commit
* gitSpell: executes an arbitrary git command
* moveSpell: moves/renames a file in the git working directory
* mergeSpell: merge two branches
* pushSpell: push to the remote repository
* removeAndCommitSpell: removes files and commit the change

## Symbols

The spells are represented by symbols in the gitalchemist file:

* symbolInit: "init\_bare\_repo"
* symbolCreateFile: "create\_file"
* symbolAdd: "add"
* symbolCommit: "commit"
* symbolGit: "git"
* symbolCreateAddCommit: "create\_add\_commit"
* symbolMerge: "merge"
* symbolPush: "push"
* symbolMove: "mv"
* symbolRemoveCommit: "remove\_and\_commit"


# laboratory.go

The laboratory file contains some common settings used in different places.


# Transmute

Transmute is the process of doing alchemy. It needs:

* a formula
* the options
* a logger to report the work

It returns an error if something went wrong.

# assistant

assistant is an interface that defines some 'low-level' methods of things
that have to be done:

* copying a file
* creating a directory
* executing a git command
* writing messages to the log

There are three implementations of assistants


## adept

An adept is a skilled assistant that can execute the instructions it gets.
If something fails, the adept reports an error.

The adept is used for running gitAlchemist in normal mode.

The adept forwards the logging task to the novice.


## novice

An novice is an unskilled assistant that does not execute the instructions,
it just reports them. It never returns an error.

The novice is used for running gitAlchemist in test mode.
 

## assistantSpy

An assistantSpy is a test double that records the calls, but does not
execute them. It also does no logging.
It can return an error on demand.

The assistantSpy is used for the unit tests.


# magic books

Magic spells are recorded in magic books.

## Read

Read is a function to read a single formula from a gitalchemist yaml file.
It returns the formula and an error if something went wrong.

## ListPages

ListPages is a function that returns a list of gitalchemist.yaml files 
for the provided task list that are found in the configuration directories

## ListBookContent

ListBookContent is a function that returns a list of all
gitalchemist.yaml files that are found  in the configuration directory.

# Options

Options control the behavior of the alchemy transmutation.

* TaskDir: directory where the files for the current are located
* RepoDir: directory where the bare repository is located
* CfgDir: directory of the configuration definitions
* Verbose: verbose logging (including debug messages)
* Test: run in test mode (novice)
* ExecuteSpells: execute only the first # spells (1-based)
* taskName: the name of the task to execute
* cloneTo: directory of the repository clone (set by initRepoSpell)
* numberOfSpells: number of spells (from Formula, set in Transmute)
* currentSpell: number of the current step (1-based)  (set in Transmute)


# Errors

There are several kinds of errors that can be returned.

* YamlDecodeError: the yaml structure is not valid
* MissingValueError: some ingredients are missing
* InvalidValueError: some ingredients are not usable
* ExecError: the execution of a spell failed
* IOError: an low-level i/o operation failed

System errors are wrapped into qualified errors.
They are not considered during unit tests to achieve operating system 
independent development.

Basic and system errors can be unwrapped if necessary.


# mortalLogger

The mortalLogger records the spells so humble mortals can admire what happened.

It provides infos in default mode and debug messages in verbose mode.


# check.ErrorString

Function "ErrorString" in package "check" is a test helper function that allows simple
error checking using the retrieved error and the wanted (expected) error message. 
If the wanted error message is the empty string, it is expected that
the error is nil.

It skips the rest of the test function if an error is wanted.

In cases where we want to test the other results even in error cases, the call
to check.ErrorString is placed at the end of each subtest.

For more details, see [README.md](../check/README.md).

# Testing

## unit tests

Executed with one of these calls:

* make test: run all unit tests
* make testv: run all unit tests in verbose mode
* make cover: run test coverage analysis, results in cover.html

Note that no real git command is executed in the unit tests!
Git commands are executed in the accpetance tests in the 
[main](../../cmd/gitalchemist/) package.

## testing frameworks

All tests use the standard go test framework: the go test tool and the
"testing" package from the standard library.

Additionally, it uses the Diff function of the 
[cmp](github.com/google/go-cmp/cmp) package provided by google
and a custom function to compare error messages (see package [check](../check/)).

Additional testing frameworks are not necessary, they just add dependencies
to the project.


## run unit tests

Example: run all unit tests

    pkg/alchemist> make test
    go test
    ok
    PASS
    ok      github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist 0.023s

Example: run a specific subtest in verbose mode:

    pkg/alchemist> go test -v -run TestListBookContent/dir_not_found
    === RUN   TestListBookContent
    === RUN   TestListBookContent/dir_not_found
        book_test.go:93: INFO: got error: read dir: open not found: no such file or directory
        book_test.go:93: skip subsequent statements due to wanted error
    --- PASS: TestListBookContent (0.00s)
        --- SKIP: TestListBookContent/dir_not_found (0.00s)
    PASS
    ok      github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist 0.003s


## test coverage

Example: examine code coverage

	go test -cover -coverprofile cover.txt 
	ok
	PASS
	coverage: 98.2% of statements
	ok  	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist	0.027s
	go tool cover -func=cover.txt
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/adept.go:25:			newAdept	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/adept.go:34:			git		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/adept.go:56:			makedir		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/adept.go:68:			copy		85.7%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/adept.go:115:			copyFile	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/adept.go:166:			examine		94.4%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/adept.go:204:			examineTarget	90.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/book.go:16:			ListPages	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/book.go:39:			ListBookContent	90.9%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:9:			Error		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:19:			Error		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:31:			Error		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:36:			Unwrap		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:51:			Error		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:67:			Unwrap		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:81:			Error		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/error.go:86:			Unwrap		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/formula.go:20:			Transmute	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/formula.go:74:			UnmarshalYAML	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/formula.go:132:			unmarshalCaster	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/formula.go:145:			Read		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/formula.go:157:			readFormula	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/laboratory.go:26:		getAuthor	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/logger.go:12:			debug		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/logger.go:24:			info		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/novice.go:18:			newNovice	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/novice.go:29:			git		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/novice.go:36:			copy		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/novice.go:43:			makedir		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spelladd.go:11:			validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spelladd.go:19:			cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellcommit.go:12:		validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellcommit.go:25:		cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellcreateaddcommit.go:19:	validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellcreateaddcommit.go:48:	cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellcreatefile.go:14:		validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellcreatefile.go:25:		cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellgit.go:14:			validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellgit.go:22:			cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellgit.go:38:			splitArgs	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellinitrepo.go:15:		validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellinitrepo.go:26:		cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellmerge.go:16:		validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellmerge.go:27:		cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellmove.go:14:		validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellmove.go:25:		cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellpush.go:11:		validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellpush.go:16:		cast		100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellremoveandcommit.go:17:	validate	100.0%
	github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/spellremoveandcommit.go:32:	cast		100.0%
	total:													(statements)	98.2%
	go tool cover -html=cover.txt -o cover.html
	echo "chromium-browser cover.html &"
	chromium-browser cover.html &
	echo "firefox cover.html &"
	firefox cover.html &


The file cover.html provides the code coverage analysis by each go file. 
You can chose the file from a dropdown list.

![../../doc/images/example_testcoverage_filelist.png](../../doc/images/example_testcoverage_filelist.png)

Code paths covered by test cases are shown in green, paths not tested are shown in red.
Not executable code parts (like declarations) are shown in grey.

![../../doc/images/example_testcoverage.png](../../doc/images/example_testcoverage.png)

The code coverage run currently does not exeute the docker tests.


## docker tests

To provide some special system conditions in a stable manner, 
tests can run in docker containers.

* make dockertest: run docker tests
* make dockertestv: run docker tests verbose

The docker test are executed with go test by using build tags.

make dockertest executes the command: 

    go test -tags teststartdocker -run TestStartDocker

The test that is started directly is **TestStartDocker**, it has the
build tag **teststartdocker** so it is not compiled in standard unit testing.

TestStartDocker starts and configures a docker container and executes 
another go test call in the container. It calls all tests that
start with **TestDocker** with the built tag **testrundocker**.

    go test -tags testrundocker -run TestDocker

If the tests ran successfully, the docker container is stopped
except when runninng go test in verbose mode.
If the tests failed, the container is not stopped and can be examined.

If it was not stopped, it runs for 5 minutes and then it terminates itself.

The command to login to the container is displayed on the terminal
as a go test INFO message.

Example:

    pkg/alchemist> make dockertestv 
    go test -v -tags teststartdocker -run TestStartDocker
    === RUN   TestStartDocker
        dockerstart_test.go:37: INFO: executing docker [exec -u testuser 04b8e02ac9ad14390b37c4709673e83b5f400f3a87cb9cbca336f39f3caad30f go test -tags testrundocker -run TestDocker -v]
        dockerstart_test.go:45: INFO: docker exec -it 04b8e02ac9ad14390b37c4709673e83b5f400f3a87cb9cbca336f39f3caad30f /bin/bash
    === RUN   TestDockerAdeptFileCopy
    === RUN   TestDockerAdeptFileCopy/copy_target_not_writeable
        docker_test.go:27: INFO: test preparation: remove : no such file or directory
        docker_test.go:27: INFO: got error: create target file: open /target.txt: permission denied
        docker_test.go:27: skip subsequent statements due to wanted error
    --- PASS: TestDockerAdeptFileCopy (0.00s)
        --- SKIP: TestDockerAdeptFileCopy/copy_target_not_writeable (0.00s)
    PASS
    ok      github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist 0.004s
    --- PASS: TestStartDocker (8.70s)
    PASS
    ok      github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist 8.702s

To examine the conainer, just execute the docker exec command that is displayed
in the test output:

    docker exec -it 04b8e02ac9ad14390b37c4709673e83b5f400f3a87cb9cbca336f39f3caad30f /bin/bash


# go doc

## exported types and functions

The standard usage of this package is to Read a gitalchemist.yaml file
to retriev a formula and to call Transmute with the formula, the options
and a logger.


	package alchemist // import "github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist"

	Package alchemist contains all the elements to do alchemistry.

	CONSTANTS

	const FormulaFileName = "gitalchemist.yaml"
		 FormulaFileName defines the name of gitalchemist formula files.


	FUNCTIONS

	func ListBookContent(cfgdir string) ([]string, error)
		 ListBookContent lists the gitalchemy.yaml files of all subdirectories of the
		 cfgdir directory.

		 It returns an error if it can not access the directories.

		 It does not return an error if a directory does not contain a
		 gitalchemy.yaml file, but this directory is not included in the returned
		 list.

	func ListPages(cfgdir string, tasklist ...string) ([]string, error)
		 ListPages lists the formula file for the provided tasks. The tasklist
		 is a task (directory) name or a comma separated list of task names. Each
		 directory must contain a gitalchemy.yaml file.

		 It returns an error if the info of the gitalchemy.yaml file can not be
		 accessed.

	func Transmute(f Formula, opt Options, logger *log.Logger) error
		 Transmute creates the git repositories according to the formula and the
		 options. Messages are writte to the logger.


	TYPES

	type ExecError struct {
		Cmd  string
		Args []string
		Err  error
	}
		 ExecError signals an error that happend during execution. It keeps the
		 underlying error that was returned by the os package.

	func (e ExecError) Error() string
		 Error returns the command, the arguments and the message of the underlying
		 error. It implents the error interface.

	func (e ExecError) Unwrap() error
		 Unwrap returns the underlying error.

	type Formula struct {
		Title    string  `yaml:"title"`
		Commands symbols `yaml:"commands"`
	}
		 Formula contains the instructions from the gitalchemy.yaml file.

	func Read(filename string) (Formula, error)
		 Read reads the file with the definitions and returns the content as a
		 Formula object. If an error occurs, it is returned.

	type IOError struct {
		Cmd string
		Arg string
		Err error
	}
		 IOError signals an error that happend during i/o operations. It keeps the
		 underlying error that was returned by the os package.

	func (e IOError) Error() string
		 Error returns the command and the message from the os which usually contains
		 all relevant info.

	func (e IOError) Unwrap() error
		 Unwrap returns the underlying error.

	type InvalidValueError struct {
		Variable, Reason string
	}
		 InvalidValueError signals an invalid value in the recipe definition.

	func (e InvalidValueError) Error() string
		 Error implents the error interface.

	type MissingValueError string
		 MissingValueError signals a missing value in the recipe definition.

	func (e MissingValueError) Error() string
		 Error implents the error interface.

	type Options struct {
		TaskDir       string // directory where the definition is located
		RepoDir       string // directory where the bare repository is located
		CfgDir        string // directory of the configuration definitions
		Verbose       bool   // verbose logging
		Test          bool   // test mode
		ExecuteSpells int    // execute only the first # steps

		// Has unexported fields.
	}
		 Options defines the options provided to the execution of the commands.

	type YamlDecodeError struct {
		Element string
		Err     error
	}
		 YamlDecodeError signals an error during yaml decoding. It keeps the
		 underlying error that was returned by the yaml package.

	func (e YamlDecodeError) Error() string
		 Error implents the error interface.

	func (e YamlDecodeError) Unwrap() error
		 Unwrap returns the underlying error.




## all types and functions

This listing contains both exported and unexported types and functions.

	package alchemist // import "github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist"

	Package alchemist contains all the elements to do alchemistry.

	CONSTANTS

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
		 A Formula file can contains these symbols for spells.

	const (
		linuxGitCmd   = "git"
		windowsGitCmd = "git.exe"
	)
		 git commands on linux and windows.

	const FormulaFileName = "gitalchemist.yaml"
		 FormulaFileName defines the name of gitalchemist formula files.

	const defaultBranch = "main"
		 defaultBranch is the git default branch name.

	const defaultUser = "red"
		 defaultUser is used when a new repo is initialzeed.

	const dirMode = 0755
		 dirMode defines the access mode for new directories.

	const gitCmd = linuxGitCmd
		 use the linux git command unless we compile for windows.

	const gitCommitDateFormat = "format:relative:5.hours.ago"

	VARIABLES

	var author = map[string]string{
		"red":       "Richard Red",
		"blue":      "Betty Blue",
		"green":     "Garry Green",
		"api":       "Alissa Api",
		"blacklist": "Benjamin Blacklist",
		"config":    "Carry Config",
	}
		 authors provides dummy author names for git.

	var email = map[string]string{
		"red":       "richard@pw-compa.ny",
		"blue":      "betty@pw-compa.ny",
		"green":     "garry@pw-compa.ny",
		"api":       "api@pw-compa.ny",
		"blacklist": "blacklist@pw-compa.ny",
		"config":    "config@pw-compa.ny",
	}
		 emails provides dummy email addresses for git.

	var regexpSplitCreateAddCommit = regexp.MustCompile(`\s*=>\s*`)
		 regexpSplitCreateAddCommit defines the regular expression for splitting the
		 source and target speciffication in the files list.


	FUNCTIONS

	func ListBookContent(cfgdir string) ([]string, error)
		 ListBookContent lists the gitalchemy.yaml files of all subdirectories of the
		 cfgdir directory.

		 It returns an error if it can not access the directories.

		 It does not return an error if a directory does not contain a
		 gitalchemy.yaml file, but this directory is not included in the returned
		 list.

	func ListPages(cfgdir string, tasklist ...string) ([]string, error)
		 ListPages lists the formula file for the provided tasks. The tasklist
		 is a task (directory) name or a comma separated list of task names. Each
		 directory must contain a gitalchemy.yaml file.

		 It returns an error if the info of the gitalchemy.yaml file can not be
		 accessed.

	func Transmute(f Formula, opt Options, logger *log.Logger) error
		 Transmute creates the git repositories according to the formula and the
		 options. Messages are writte to the logger.

	func examine(from, to string) (target, dir string, recursive bool, err error)
		 examine checks the source ("from") and target ("to") of a copy call.
		 It returns:

			- the new target for the copy
			- a directory name if one should be created (might be empty)
			- a bool to indicate recursive copy
			- an error if one occurred

		 The source is determined with the following logic.

			- If from does not exist, it returns an error.
			- If from is a file, then recursive will be false.
			- If from is a directory, then recursive will be true.

		 The target is determined with the following logic.

		 If the target exists:

			- If the check of "to" fails, the error is returned.
			- If to is a file, it is returned
			- If to is directory, it appends the last element from 'from' to 'to'

		 If the target does not exist:

			- If the source is a directory or if the source is a file and to ends with
			  a slash, it is used a a directory
			- Else it is used as target file.

		 If a directory to the target should be created, dir contains the name.

	func examineTarget(to string) (exist, isdir bool, err error)
		 examineTarget examines the to element. It returns if it does exist and if it
		 is or should be a directory.

	func getAuthor(name string) string
		 getAuthor tries to get the author with mail address. if is not known,
		 the name is returned as is.

	func unmarshalCaster[T caster](node *yaml.Node) (T, error)
		 unmarshalCaster extracts a single caster from the yaml node.


	TYPES

	type ExecError struct {
		Cmd  string
		Args []string
		Err  error
	}
		 ExecError signals an error that happend during execution. It keeps the
		 underlying error that was returned by the os package.

	func (e ExecError) Error() string
		 Error returns the command, the arguments and the message of the underlying
		 error. It implents the error interface.

	func (e ExecError) Unwrap() error
		 Unwrap returns the underlying error.

	type Formula struct {
		Title    string  `yaml:"title"`
		Commands symbols `yaml:"commands"`
	}
		 Formula contains the instructions from the gitalchemy.yaml file.

	func Read(filename string) (Formula, error)
		 Read reads the file with the definitions and returns the content as a
		 Formula object. If an error occurs, it is returned.

	func readFormula(r io.Reader) (Formula, error)
		 readFormula extracts the yaml content from the provided reader.

	type IOError struct {
		Cmd string
		Arg string
		Err error
	}
		 IOError signals an error that happend during i/o operations. It keeps the
		 underlying error that was returned by the os package.

	func (e IOError) Error() string
		 Error returns the command and the message from the os which usually contains
		 all relevant info.

	func (e IOError) Unwrap() error
		 Unwrap returns the underlying error.

	type InvalidValueError struct {
		Variable, Reason string
	}
		 InvalidValueError signals an invalid value in the recipe definition.

	func (e InvalidValueError) Error() string
		 Error implents the error interface.

	type MissingValueError string
		 MissingValueError signals a missing value in the recipe definition.

	func (e MissingValueError) Error() string
		 Error implents the error interface.

	type Options struct {
		TaskDir       string // directory where the definition is located
		RepoDir       string // directory where the bare repository is located
		CfgDir        string // directory of the configuration definitions
		Verbose       bool   // verbose logging
		Test          bool   // test mode
		ExecuteSpells int    // execute only the first # steps

		// set during processing
		taskName       string // name of the task to execute
		cloneTo        string // directory of the repository clone
		numberOfSpells int    // number of steps
		currentSpell   int    // number of the current step (1-based)
	}
		 Options defines the options provided to the execution of the commands.

	type YamlDecodeError struct {
		Element string
		Err     error
	}
		 YamlDecodeError signals an error during yaml decoding. It keeps the
		 underlying error that was returned by the yaml package.

	func (e YamlDecodeError) Error() string
		 Error implents the error interface.

	func (e YamlDecodeError) Unwrap() error
		 Unwrap returns the underlying error.

	type addSpell struct {
		Files []string `yaml:"files"`
	}
		 addSpell provides adding files to the git index.

	func (s addSpell) cast(a assistant, opt Options) error
		 cast executes git add of all files in the list.

	func (s addSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type adept struct {
		novice
		exe string
	}
		 adept is an assitant that has a high skill level so it can execute all the
		 steps. The logging is delegated to the novice. The actions will be logged at
		 debug level.

		 It implements the assistant interface.

	func newAdept(l *log.Logger, opt Options) adept
		 newAdept returns an initialized adept object.

	func (a adept) copy(from, to string) error
		 copy copies the file to the target. The target directory must exist.

	func (a adept) copyFile(from, to string) error
		 copyFile copies the file to the target. The target directory must exist.

	func (m adept) debug(msg string, args ...any)
		 debug writes a debug message to the logger if verbose is true.

	func (a adept) git(dir string, args ...string) error
		 git executes the git command in the provided directory. The directory must
		 exist.

	func (m adept) info(msg string, args ...any)
		 info writes an info message to the logger.

	func (a adept) makedir(dir string) error
		 makedir creates the target dir and all missing directories on the path.

	type assistant interface {
		// git execute a git command in the provided directory
		git(dir string, args ...string) error
		// copy copies a file
		copy(from, to string) error
		// makedir creates a directory
		makedir(dir string) error

		// debug emits a debug message
		debug(msg string, args ...any)
		// info emits an info
		info(msg string, args ...any)
	}
		 assistant define the ability to execute low-level commands.

	type caster interface {
		validate() error
		cast(assistant, Options) error
	}
		 caster defines the ability to cast a spell (i.e. execute a command) from the
		 Formula. It also allows checking if the spell is valid.

	type commitSpell struct {
		Message string `yaml:"message"`
		Author  string `yaml:"author"`
	}
		 commitSpell provides committing the index to the repo.

	func (s commitSpell) cast(a assistant, opt Options) error
		 incant executes git commit.

	func (s commitSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type createAddCommitSpell struct {
		Files   []string `yaml:"files"`
		Message string   `yaml:"message"`
		Author  string   `yaml:"author"`
	}
		 createAddCommitSpell

	func (s createAddCommitSpell) cast(a assistant, opt Options) error
		 cast copies the files, adds them to the index and commits it. It uses
		 createFileSpell, addSpell and commitSpell.

	func (s createAddCommitSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type createFileSpell struct {
		Source string `yaml:"source"`
		Target string `yaml:"target"`
	}
		 createFileSpell provides copying a file to the git repo directory.

	func (s createFileSpell) cast(a assistant, opt Options) error
		 cast copies the file to the repo.

	func (s createFileSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type gitSpell struct {
		Command string `yaml:"command"`
	}
		 gitSpell provides arbitrary git commands.

	func (s gitSpell) cast(a assistant, opt Options) error
		 cast executes an arbitrary git command.

	func (s gitSpell) splitArgs() []string
		 splitArgs splits args like the shell: double quoted text can contain spaces.

	func (s gitSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type initRepoSpell struct {
		Bare    string `yaml:"bare"`
		CloneTo string `yaml:"clone_to"`
	}
		 initRepoSpell provides initializing a bare repo and cloning it.

	func (s initRepoSpell) cast(a assistant, opt Options) error
		 cast creates a bare repo and clones and initializes it.

	func (s initRepoSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type mergeSpell struct {
		Source       string `yaml:"source"`
		Target       string `yaml:"target"`
		DeleteSource bool   `yaml:"delete_source"`
	}
		 mergeSpell provides merging two git branches.

	func (s mergeSpell) cast(a assistant, opt Options) error
		 cast executes a git merge.

	func (s mergeSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type mortalLogger struct {
		*log.Logger
		verbose bool
	}
		 mortalLogger provides human-readable logging.

	func (m mortalLogger) debug(msg string, args ...any)
		 debug writes a debug message to the logger if verbose is true.

	func (m mortalLogger) info(msg string, args ...any)
		 info writes an info message to the logger.

	type moveSpell struct {
		Source string `yaml:"source"`
		Target string `yaml:"target"`
	}
		 moveSpell provides moving a file in the git clone.

	func (s moveSpell) cast(a assistant, opt Options) error
		 cast moves a file in the git working directory.

	func (s moveSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type novice struct {
		mortalLogger
	}
		 novice is an assistant that does not execute the commands, it just reports
		 the steps that should be executed (on debug level).

		 It implements the assistant interface.

		 It is used for test mode. All methods return a nil error.

	func newNovice(l *log.Logger, opt Options) novice
		 newNovice returns an initialized novice object.

	func (n novice) copy(from, to string) error
		 copy emits a debug message with the parameters. It implements the assistant
		 interface.

	func (m novice) debug(msg string, args ...any)
		 debug writes a debug message to the logger if verbose is true.

	func (n novice) git(dir string, args ...string) error
		 git emits a debug message with the parameters. It implements the assistant
		 interface.

	func (m novice) info(msg string, args ...any)
		 info writes an info message to the logger.

	func (n novice) makedir(dir string) error
		 makedir emits a debug message with the parameters. It implements the
		 assistant interface.

	type pushSpell struct {
		Main bool `yaml:"main"`
	}
		 pushSpell provides push to a remote repository.

	func (s pushSpell) cast(a assistant, opt Options) error
		 cast executes a git push if Main is true.

	func (s pushSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type removeAndCommitSpell struct {
		Files   []string `yaml:"files"`
		Message string   `yaml:"message"`
		Author  string   `yaml:"author"`
	}
		 removeAndCommitSpell provides the combined command of removing files and
		 commiting them.

	func (s removeAndCommitSpell) cast(a assistant, opt Options) error
		 cast calls git rm to all the files and commits the result.

	func (s removeAndCommitSpell) validate() error
		 validate checks the values and reports an error if something is missing.

	type spellHint struct {
		dir  string
		args []string
	}
		 spellHint provides the information for casting some spells. It can be used
		 to create lists of arguments and process them in a loop.

	type symbols struct {
		cloneTo string
		spells  []caster
	}
		 symbols contains a list of spells.

	func (c *symbols) UnmarshalYAML(value *yaml.Node) (err error)
		 UnmarshalYAML extracts the commands from the yaml definition into a list of
		 casters. It also checks for completeness and correctness.


