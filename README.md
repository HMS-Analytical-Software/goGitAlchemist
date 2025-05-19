# go GitAlchemist

GitAlchemist is a framework for creating git repositories from config 
files on local disk with a predefined history of commits. 
The primary purpose of this tool is setting up interactive tasks for 
git workshops or git tutorials.

## Requirements

go GitAlchemist compiles to binaries, so it only needs:

* a gitalchemist binary for the current operating system
* a git installation
* and the gitalchemy.yaml definition files and all the repository content created by you

## Usage

Gitalchemist has to be invoked with either

* one or more task names 
* -runall: execute all tasks in the configuration directory
* -clean: remove the target directory

A task or *spell* is a directory that contains a gitalchemist.yaml definition
file and all the files that are used in the definition.
You can find examples in the directory 
[cmd/gitalchemist/testdata](cmd/gitalchemist/testdata).

Options:

* -test: run in test mode, do not execute commands, just show them
* -verbose: run in verbose mode
* -targetdir: write the git repos to this directory
* -cfgdir: search tasks here

The cfgdir is prefixed to the task names, so these calls are equivalent:

```bash
./gitalchemist -cfgdir testdata basic_workflow cmd_merge
./gitalchemist testdata/basic_workflow testdata/cmd_merge
```

Online help:

```
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
```

## Exit codes

The exit code of the program is determined by the kind of error that happened:

* 0: no errors
* 1: invalid command line parameters
* 2: missing values in gitalchemist file
* 3: invalid values in gitalchemist file
* 4: invalid yaml in gitalchemist file
* 5: command execution error
* 6: i/o error 
* 42: other


## Features

Currently, the following commands are supported:

* **init\_bare\_repo**: create a bare repo and clone it
* **create\_file**: copy a file to the working directory
* **add**: add files to the index
* **commit**: commit the index
* **create\_add\_commit**: combined create\_file, add, and commit
* **git**: execute arbitrary git command
* **merge**: merge two branches
* **push**: push to the remote repo
* **mv**: move a file within git
* **remove\_and\_commit**: remove files and commit change


## Example: gitalchemist.yaml

The structure of the gitalchemist.yaml file is defined like this.
All elements of a command are reqired unless the ones that are
marked as 'optional'.

```yaml
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
```

## Call example

Call of one of the examples in cmd/gitalchemist/testdata:

```
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
```

Call in verbose mode:

```
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
```

Clean created git repos:

```
cmd/gitalchemist > ./gitalchemist -clean
2025/02/23 10:36:11 [INFO] remove cwd
2025/02/23 10:36:11 [INFO] ok
```

## Installation

It is available for the operating systems and processor architectures
supported by go, e.g Linux, Windows and MacOsX.

Installation is not yet developed, for the time being exectutables are
provided as GitHub releases.

## Development

see [DEVELOP.md](DEVELOP.md)

## Python Implementation

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

There are two major differences:

### Git commands

Go GitAlchemist does not allow execution of arbitrary commands.
When using the 'git' command definition, the operating system
specific executable ('git' or 'git.exe') is executed.

To be backwards compatible, it allows that the command defined in the 
gitalchemist.yaml file starts with 'git' (or 'git.exe'). This first
argument is removed during processing.

### Timestamp directory

Go GitAlchemist does not use a timestamp directory. Each created git repo
is placed within the target directory, by default "./cwd", in a directory
that corresponds to the Title definition in the gitalchemist file.

This implies that the created directories must be cleaned up before a rerun.
It avoids cluttering the directory space.

To remove all temporary git repos, use the -clean command line switch.

