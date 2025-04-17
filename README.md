# go GitAlchemist

GitAlchemist is a small framework for creating git repositories from config 
files on local disc with a predefined history of commits. 
The primary purpose of this tool is setting up interactive tasks for 
git workshops or git tutorials.

This repository contains an implementation of 
[GitAlchemist](https://github.com/HMS-Analytical-Software/GitAlchemist)
in go.

Go GitAlchemist compiles to binaries, so it does not need a python installation 
at run time.

Go GitAlchemist tries to be compatible with GitAlchemist by using
the same gitalchemist.yaml file format that defines
a 'formula' for git alchemistry. A valid GitAlchemist file should
be usable with go GitAlchemist without changes.

The internal code structure, the log output and the command line parameters 
however are different.

# Differences

There are three major differences:

## git commands

Go GitAlchemist does not allow execution of arbitrary commands.
When using the 'git' command definition, the operating system
specific executable ('git' or 'git.exe') is executed.

To be backwards compatible, it allows that the command defined in the 
gitalchemist.yaml file starts with 'git' (or 'git.exe'). This first
argument is removed during processing.

## timestamp directory

Go GitAlchemist does not use a timestamp directory. Each created git repo
is placed within the target directory, by default "./cwd", in a directory
that corresponds to the Title definition.

This implies that the created directories must be cleaned up befor a rerun.
It avoids cluttering the directory space.

To remove all temporary created git repos, use the -clean command line switch.

## requirements

go GitAlchemist compiles to binaries, so it needs only:

* a gitalchemist binary for the current operating system
* a git installation
* and the gitalchemy.yaml definition files and all the repository content created by you


## Usage

gitalchemist has to be called with either

* one or more task names 
* -runall: execut all tasks in cfgdir
* -clean: remove targetdir

A task is a directory that contains a gitalchemist.yaml definition
file and all the files that are used in the definition.
You can find examples in the directory 
[cmd/gitalchemist/testdata](cmd/gitalchemist/testdata).

Options:

* -test: run in test mode, do not execute commands, just show them
* -verbose: run in verbose mode
* -targetdir: write git repos to this directory
* -cfgdir: search tasks here

The cfgdir is prefixed to the task names, so these calls are equivalent:

    ./gitalchemist -cfgdir testdata basic_workflow cmd_merge
    ./gitalchemist testdata/basic_workflow testdata/cmd_merge

online help:

    cmd/gitalchemist> ./gitalchemist -h
    usage: ./gitalchemist <path/to/dir> [ path/to/dir ... ]]
    usage: ./gitalchemist -runall
    usage: ./gitalchemist -clean

    The directories must contain a definition file named "gitalchemist.yaml" 
    and all the files that are used in the definition.

    usage of ./gitalchemist:
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

## exit codes

The exit code of the program is determined by the kind of error that happened:

* 0: no errors
* 1: invalid command line parameters
* 2: missing values in gitalchemist file
* 3: invalid values in gitalchemist file
* 4: invalid yaml in gitalchemist file
* 5: command execution error
* 6: i/o error 
* 42: other

## features

Currently, the following commands are supported:

* init\_bare\_repo: create a bare repo and clone it
* create\_file: copy a file to the working directory
* add: add files to the index
* commit: commit the index
* create\_add\_commit: combined create\_file, add, and commit
* git: execute arbitrary git command
* merge: merge two branches
* push: push to the remote repo
* mv: move a file within git
* remove\_and\_commit: remove files and commit change


## Example: gitalchemist.yaml

The structure of the gitalchemist.yaml file is defined like this.
All elements of a command are reqired unless the ones that are
marked as 'optional'.


    title: test_workflow # must match directory name!
    commands:
      - init_bare_repo:
          bare: remotes/create_add_commit
          clone_to: workflow
      - create_file:
          source: files/project_plan_v1.md
          target: project_plan.md
      - add:
          files: 
          - project_plan.md
      - commit:
          message: Added first file
          author: red
      - create_add_commit:
          files:
          - files/project_plan_v3.md => project_plan.md
          message: removed unnecessary parts of the project plan
          author: red
      - create_add_commit:
          files:
          - files/folder1 => folder1/
          message: added folder1
          author: red
      - git:
          command: "commit -m \"my message\""
      - git:
          # legacy mode with explicit git
          command: "git commit -m \"my message\"" 
      - merge:
          source: feature/start_project
          target: main
          # optional, defaults to false if missing
          delete_source: true 
      - push:
          main: true
      - mv:
          source: main.py
          target: generator.py
      - remove_and_commit:
          files:
          - notes-timeline.txt
          message: clean up timeline notes
          author: red

## call example


Call of one of the examples in cmd/gitalchemist/testdata:

    cmd/gitalchemist > ./gitalchemist -cfgdir testdata basic_workflow 
    2025/02/23 10:34:55 [INFO] execute formula basic_workflow
    2025/02/23 10:34:55 [INFO] 1/4: make directory cwd/remotes/basic_workflow
    2025/02/23 10:34:55 [INFO] 1/4: init --bare --initial-branch=main .
    2025/02/23 10:34:55 [INFO] 1/4: clone remotes/basic_workflow basic_workflow
    2025/02/23 10:34:55 [INFO] 1/4: remote set-url origin ../remotes/basic_workflow
    2025/02/23 10:34:55 [INFO] 1/4: config user.name Richard Red
    2025/02/23 10:34:55 [INFO] 1/4: config user.email richard@pw-compa.ny
    2025/02/23 10:34:55 [INFO] 1/4: config init.defaultBranch main
    2025/02/23 10:34:55 [INFO] 2/4: copy testdata/basic_workflow/files/project_plan_v1.md to cwd/basic_workflow/project_plan.md
    2025/02/23 10:34:55 [INFO] 3/4: add 1 files
    2025/02/23 10:34:55 [INFO] 4/4: commit
    2025/02/23 10:34:55 [INFO] ok

Call in verbose mode:

    cmd/gitalchemist> ./gitalchemist -verbose -cfgdir testdata basic_workflow 
    2025/02/23 10:37:16 [INFO] execute formula basic_workflow
    2025/02/23 10:37:16 [INFO] 1/4: make directory cwd/remotes/basic_workflow
    2025/02/23 10:37:16 [DEBUG] makedir "cwd/remotes/basic_workflow"
    2025/02/23 10:37:16 [INFO] 1/4: init --bare --initial-branch=main .
    2025/02/23 10:37:16 [DEBUG] "cwd/remotes/basic_workflow": git []string{"init", "--bare", "--initial-branch=main", "."}
    2025/02/23 10:37:16 [DEBUG] Initialized empty Git repository in /home/jochen/work/hms/czemmel-goGitAlchemist/cmd/gitalchemist/cwd/remotes/basic_workflow/
    2025/02/23 10:37:16 [INFO] 1/4: clone remotes/basic_workflow basic_workflow
    2025/02/23 10:37:16 [DEBUG] "cwd": git []string{"clone", "remotes/basic_workflow", "basic_workflow"}
    2025/02/23 10:37:16 [DEBUG] Cloning into 'basic_workflow'...
    2025/02/23 10:37:16 [DEBUG] warning: You appear to have cloned an empty repository.
    2025/02/23 10:37:16 [DEBUG] done.
    2025/02/23 10:37:16 [INFO] 1/4: remote set-url origin ../remotes/basic_workflow
    2025/02/23 10:37:16 [DEBUG] "cwd/basic_workflow": git []string{"remote", "set-url", "origin", "../remotes/basic_workflow"}
    2025/02/23 10:37:16 [INFO] 1/4: config user.name Richard Red
    2025/02/23 10:37:16 [DEBUG] "cwd/basic_workflow": git []string{"config", "user.name", "Richard Red"}
    2025/02/23 10:37:16 [INFO] 1/4: config user.email richard@pw-compa.ny
    2025/02/23 10:37:16 [DEBUG] "cwd/basic_workflow": git []string{"config", "user.email", "richard@pw-compa.ny"}
    2025/02/23 10:37:16 [INFO] 1/4: config init.defaultBranch main
    2025/02/23 10:37:16 [DEBUG] "cwd/basic_workflow": git []string{"config", "init.defaultBranch", "main"}
    2025/02/23 10:37:16 [INFO] 2/4: copy testdata/basic_workflow/files/project_plan_v1.md to cwd/basic_workflow/project_plan.md
    2025/02/23 10:37:16 [DEBUG] makedir "cwd/basic_workflow"
    2025/02/23 10:37:16 [DEBUG] copy "testdata/basic_workflow/files/project_plan_v1.md" to "cwd/basic_workflow/project_plan.md"
    2025/02/23 10:37:16 [DEBUG] makedir "cwd/basic_workflow"
    2025/02/23 10:37:16 [INFO] 3/4: add 1 files
    2025/02/23 10:37:16 [DEBUG] "cwd/basic_workflow": git []string{"add", "project_plan.md"}
    2025/02/23 10:37:16 [INFO] 4/4: commit
    2025/02/23 10:37:16 [DEBUG] "cwd/basic_workflow": git []string{"commit", "--date=format:relative:5.hours.ago", "-m", "Added first file", "--author=Richard Red <richard@pw-compa.ny>"}
    2025/02/23 10:37:16 [DEBUG] [main (root-commit) 92b656e] Added first file
    2025/02/23 10:37:16 [DEBUG]  Date: Sun Feb 23 05:37:16 2025 +0100
    2025/02/23 10:37:16 [DEBUG]  1 file changed, 75 insertions(+)
    2025/02/23 10:37:16 [DEBUG]  create mode 100644 project_plan.md
    2025/02/23 10:37:16 [INFO] ok


Clean created git repos:

    cmd/gitalchemist > ./gitalchemist -clean
    2025/02/23 10:36:11 [INFO] remove cwd
    2025/02/23 10:36:11 [INFO] ok


# Installation

It is available for the operating systems and processor architectures
supported by go, e.g Linux, Windows and MacOsX.

Installation is not yet developed, for the time being exectutables are
provided as GitHub releases.


# Development

* documentation for the main package: see [README.md](./cmd/gitalchemist/README.md)
* documentation for package "alchemist": see [README.md](./pkg/alchemist/README.md)
* documentation for package "check": see [README.md](./pkg/check/README.md)


# tools

Required for development

* go compiler toolkit ( >= go1.24.0): test and build software
* make: execute commands
* govulncheck: check for security vulnerabilities
* golangci-lint: execute static code checks
* git: execute the acceptance tests, and of course ;-) manage code
* docker: execute some special unit tests (optional)
* godoc: display code documentation in browser (optional)

The usage of a go specific editor plugins is recommended:

* vs code: go extension
* vim: vim-go plugin

The plugins should be configured to format the code and organize the 
imports automatically when a file is saved.

Improper formatted code is not acceptable.

The creation of releases currently requires also the git and GitHub cli ('gh').

# Unit tests

Package level unit tests can be called in the package
or command directories:

* make test: execute all unit tests
* make testv: verbose execution
* go test -run TestXXX: run only TestXXX test
* go test -run TestXXX/subtestYYY: run only TestXXX/subtestYYY test


# docker tests

Some system conditions are hard to test on a real system in a portable
manner. Some of these tests can be executed in a docker container.

* make dockertest: run specific unit tests in docker
* make dockertestv: verbose mode

For more details, see [README.md](./pkg/alchemist/README.md)

# Acceptance tests

The tests with the compiled binary are ecuted in the
directory cmd/gitalchemist, 

* make acctest TASK=cmd\_mv: run cmd\_mv test (default is basic\_workflow)
* make acctestvh TASK=cmd\_mv: run one test in verbose mode
* make acctestall: run all test tasks
* make acctestallv: run all test tasks in verbose mode
* make acctestrunall: run with runall flag
* make acctestrunallv: run with runall flag in verbose mode

For more details, see [README.md](./cmd/gitalchemist/README.md)

# Test coverage

Test coverage of unit tests can be checked in the respective
pakage directory:

* make cover

Test coverage creates two files:

* cover.txt: intermediate results
* cover.html:  html test coverage report

These files are ignored by git.

The test coverage report can be viewed directly with a web browser, e.g.:

*  firefox cover.html &

For more details, see [README.md](./cmd/gitalchemist/README.md)


# Static tests and module unit tests

Some test can be executed in the main directory:

* make vulncheck: check for vulnerable (insecure) code using govulncheck
* make vet: static code analysis using go vet
* make lint: static code analysis using golanci-lint
* make test: run all unit tests
* make testall: run vulncheck, vet, lint and test

make test only executes the standard unit tests, not the docker
or acceptance tests.

## example output

    > make testall
    govulncheck ./...
    No vulnerabilities found.
    go vet -lostcancel=false ./...
    golangci-lint run
    go test ./...
    ok      github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/cmd/gitalchemist      0.003s
    ok      github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist 0.021s
    ok      github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/check     0.003s

# github workflows

There are two github workflows that execute tests.

## go-unit-test

The github workflow go-unit-test executes the go unit tests for the 
complete module in verbose mode. 

It is executed each time when a push of go files takes place.
It can also be triggered on demand.

The tests run on Linux, windows, and macosx (darwin).

## gitalchemist-acceptance-tests

The github workflow gitalchemist-acceptance-tests executes a complete
acceptanct test. It can be triggered on demand using the GitHub
website or by using the gh command line client.

The tests can run on Linux, windows, or macosx (darwin).

For more details, see [README.md](./cmd/gitalchemist/README.md)

# Testdata files

In go, test files are placed in directories named 'testdata'.
The go tools know about the special meaning of this name.

The testdata directories contain two types of files or subdirectories:

* static: they are tracked in git and used readonly in the tests.
* dynamic: created or modified during tests, ignored by git.


# Testing without make

The makefile content is pretty simple. The commands can be
executed directly or with a script if make is not available.

# View the code documentation in browser

* make docserver: starts the godoc server on localhost:6060 and starts a browser

The code of each package can be viewed on a single page.
Per default, only exported objects are displayed.
If you append "?m=all" to the url, also unexported elements are shown.

## Example: start page

http://localhost:6060/pkg/github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/

![doc/images/example\_godoc\_mainpage.png](doc/images/example_godoc_mainpage.png)

## Example: exported only

http://localhost:6060/pkg/github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/

![doc/images/example\_godoc\_exported.png](doc/images/example_godoc_exported.png)

## Example: with unexported

http://localhost:6060/pkg/github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/pkg/alchemist/?m=all

![doc/images/example\_godoc\_unexported.png](doc/images/example_godoc_unexported.png)

# Operating system specific code and tests

Some code is specific for certain operating systems (e.g. git vs. git.exe).
This is handled using go build tags or special file names.

* //go:build !windows: if placed at the top of a go file: do not compile on windows
* xx\_windows.go: only compiles for windows.

It is developed on Linux and currently not yet tested on windows.

# Software BOM

## go binaries

go provides some metadata in every binary by default.
The information can be displayed by calling go version.

    cmd/gitalchemist> go version -m gitalchemist
    gitalchemist: go1.24.0
            path    github.com/HMS-Analytical-Software/czemmel-goGitAlchemist/cmd/gitalchemist
            mod     github.com/HMS-Analytical-Software/czemmel-goGitAlchemist       v0.1.3-0.20250301162826-66028cd15f55
            dep     gopkg.in/yaml.v3        v3.0.1  h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
            build   -buildmode=exe
            build   -compiler=gc
            build   CGO_ENABLED=1
            build   CGO_CFLAGS=
            build   CGO_CPPFLAGS=
            build   CGO_CXXFLAGS=
            build   CGO_LDFLAGS=
            build   GOARCH=amd64
            build   GOOS=linux
            build   GOAMD64=v1
            build   vcs=git
            build   vcs.revision=66028cd15f55e33da041912cdaed6b0be13904e4
            build   vcs.time=2025-03-01T16:28:26Z
            build   vcs.modified=false

This allows that the govulncheck tool is also able to scan binaries
for code vulnerabilities.

    cmd/gitalchemist> govulncheck -mode binary gitalchemist
    No vulnerabilities found.


## go modules

The code repository contains two files that contain information about
used libraries and their versions.

* go.mod: contains some module information
* go.sum: contains checksums of the dependencies

The content of the files is automatically updated by the go tools,
e.g. when downloading a module using "go install" or "go get".
Changes can be made by using 'go mod' commands.

Example:

* go mod tidy: cleans up go.mod and go.sum accoring to the current module content.

###  go.mod

Example:

    module github.com/HMS-Analytical-Software/czemmel-goGitAlchemist

    go 1.24.0

    require (
        github.com/google/go-cmp v0.6.0
        gopkg.in/yaml.v3 v3.0.1
    )

The package "github.com/google/go-cmp" is only used in the unit tests,
therefore it is not displayed in the output of the binary.

### go.sum

Example:

    github.com/google/go-cmp v0.6.0 h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=
    github.com/google/go-cmp v0.6.0/go.mod h1:17dUlkBOakJ0+DkrSSNjCkIjxS6bF9zb3elmeNGIjoY=
    gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=
    gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
    gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
    gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=

The package gopkg.in/check.v1 is a dependency of the yaml package.

