# gitalchemist

This directory contains the main package used to build
the gitalchemist executable.

It also contains unit and acceptanct tests.

# Usage

The usage of gitalchemist is documentated in the main  
[README.md](../../README.md) file.


# build

Create the binaries:

* make build: creates the executable in the local directory for the current operating system
* make linux: build binary for linux
* make win: build binary for windows
* make mac: build binary for macosx (arm processor)
* make all: build win, linux mac

# unit tests

The setting of parameters by environment variables, command line flags
and defaults can be tested.

* make test: execute unit tests
* make testv: execute unit tests in verbose mode
* make cover: test coverage

# acceptance tests

The acceptance tests test the behaviour of the compiled program.

The acceptance tests:

* remove the default output directory 'cwd'
* call go test which builds and executes the binary with the provided parameters
  and examines the created files and git repostories

The acceptance tests be called like this:

* make acctest: run basic\_workflow test (default)
* make acctest TASK=cmd\_mv: run cmd\_mv test
* make acctestv TASK=cmd\_mv: run cmd\mv test in verbose mode
* make acctestall: run all test tasks
* make acctestallv: run all test tasks in verbose mode
* make acctestrunall: run with runall flag
* make acctestrunallv: run with runall flag in verbose mode


The test configuration files are provided in the testdata directory.

Example: run test for default TASK "basic\_workflow" in verbose mode:

    cmd/gitalchemist> make acctestv 
    go build
    ./gitalchemist -clean
    2025/02/23 11:26:55 [INFO] remove cwd
    2025/02/23 11:26:55 [INFO] ok
    go test -v -tags acctest -run TestAcceptance/basic_workflow
    === RUN   TestAcceptance
    === RUN   TestAcceptance/basic_workflow
    2025/02/23 11:26:55 [INFO] execute formula basic_workflow
    2025/02/23 11:26:55 [INFO] 1/4: make directory cwd/remotes/basic_workflow
    2025/02/23 11:26:55 [INFO] 1/4: init --bare --initial-branch=main .
    2025/02/23 11:26:55 [INFO] 1/4: clone remotes/basic_workflow basic_workflow
    2025/02/23 11:26:55 [INFO] 1/4: remote set-url origin ../remotes/basic_workflow
    2025/02/23 11:26:55 [INFO] 1/4: config user.name Richard Red
    2025/02/23 11:26:55 [INFO] 1/4: config user.email richard@pw-compa.ny
    2025/02/23 11:26:55 [INFO] 1/4: config init.defaultBranch main
    2025/02/23 11:26:55 [INFO] 2/4: copy testdata/basic_workflow/files/project_plan_v1.md to cwd/basic_workflow/project_plan.md
    2025/02/23 11:26:55 [INFO] 3/4: add 1 files
    2025/02/23 11:26:55 [INFO] 4/4: commit
    2025/02/23 11:26:55 [INFO] ok
        acc_test.go:169: INFO: [rev-parse --is-bare-repository]
        acc_test.go:169: INFO: true
        acc_test.go:173: INFO: [remote -v]
        acc_test.go:173: INFO: origin       ../remotes/basic_workflow (fetch)
            origin  ../remotes/basic_workflow (push)
        acc_test.go:180: INFO: files equal: "testdata/basic_workflow/files/project_plan_v1.md" - "cwd/basic_workflow/project_plan.md"
        acc_test.go:188: INFO: [show --pretty= --name-status]
        acc_test.go:188: INFO: A    project_plan.md
        acc_test.go:188: INFO: [log --pretty=format:%s]
        acc_test.go:188: INFO: Added first file
    --- PASS: TestAcceptance (0.16s)
        --- PASS: TestAcceptance/basic_workflow (0.04s)
    PASS
    ok      github.com/HMS-Analytical-Software/goGitAlchemist/cmd/gitalchemist      0.158s

# acceptance tests in github

You can execute the acceptance tests on github by calling the action
**gitalchemist-acceptance-tests**
directly from the GitHub web site or by using the gh cli.

## Github

Select 
    Actions > gitalchemist-acceptance-tests > Run workflow
and select the operating system to use:

* all (default)
* linux
* windows
* darwin (MacOSX)


## gh cli

* make acctestworkflow: executes "make acctestallv" in a github runner using the current branch and Liunx
* make acctestworkflow TESTOS=windows: run acceptance tests for windows
* make acctestworkflow TESTOS=darwin: run acceptance tests for macosx (aka darwin)
* make acctestworkflow TESTOS=all: run acceptance tests for linux, windows, and darwin

This make call requires prior authentication of the gh client.

Example:

    cmd/gitalchemist> make acctestworkflow 
    gh workflow run gitalchemist-acceptance-tests  --ref exp/run-acceptance-tests-in-github
    ✓ Created workflow_dispatch event for gitalchemist-acceptance-tests.yml at exp/run-acceptance-tests-in-github

    To see runs for this workflow, try: gh run list --workflow=gitalchemist-acceptance-tests.yml


# release

To create a github release, just call make with the release label:

Example:

    make release REL=v0.0.1

