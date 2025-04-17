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

# code documentation

The structure of the code can be viewed on the command line by using 'go goc':

* go doc -all: shows the types, methods, functions and constants of the package.
* go doc -all -u: also shows the unexported elements


# Testing

## unit tests

Executed with one of these calls:

* make test: run all unit tests
* make testv: run all unit tests in verbose mode
* make cover: run test coverage analysis, results in cover.html

Note that no real git command is executed in the unit tests!
Git commands are executed in the accpetance tests in the 
[main](../../cmd/gitalchemist/) package.


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


