
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


# code documentation

The structure of the code can be viewed on the command line by using 'go goc'
or by starting a godoc server and view the documentation in a browser.


# Unit tests

Package level unit tests can be called in the package
or command directories.

## execute unit tests

The tests can be run using make or by executing the 
corresponding go test call:

* make test: execute all unit tests
* make testv: verbose execution
* go test -run TestXXX: run only TestXXX test
* go test -run TestXXX/subtestYYY: run only TestXXX/subtestYYY test

## testing frameworks

All tests use the standard go test framework: the go test tool and the
"testing" package from the standard library.

Additionally, they use the Diff function of the 
[cmp](https://github.com/google/go-cmp/cmp) package provided by google
and a custom function to compare error messages (see package [check](../check/)).

Some simple error comparison functions are provided in the 'check' package.


# docker tests

Some system conditions are hard to test on a real system in a portable
manner. Some of these tests can be executed in a docker container.

* make dockertest: run specific unit tests in docker
* make dockertestv: verbose mode

For more details, see [README.md](./pkg/alchemist/README.md)

# Acceptance tests

The tests with the compiled binary are ecuted in the
directory cmd/gitalchemist, 

* make acctest: run basic\_workflow test (default)
* make acctest TASK=cmd\_mv: run cmd\_mv test
* make acctestv TASK=cmd\_mv: run cmd\mv test in verbose mode
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

There are three github workflows that execute tests.

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

## govulncheck

The github workflow gitalchemist-acceptance-tests executes a
vulnerability check on the go code.

It is executed each time when a push of go files takes place.
It can also be triggered on demand.

The check is only executed on Linux.

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
* \*\_windows.go: only compiles for windows.

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

