package alchemist

import (
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// adept is an assitant that has a high skill level so it can execute
// all the steps.
// The logging is delegated to the novice.
// The actions will be logged at debug level.
//
// It implements the assistant interface.
type adept struct {
	novice
	exe string
}

// newAdept returns an initialized adept object.
func newAdept(l *log.Logger, opt Options) adept {
	return adept{
		novice: novice{mortalLogger{Logger: l, verbose: opt.Verbose}},
		exe:    gitCmd,
	}
}

// git executes the git command in the provided directory.
// The directory must exist.
func (a adept) git(dir string, args ...string) error {
	a.novice.git(dir, args...)

	cmd := exec.Command(a.exe, args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	if len(output) > 0 {
		for _, line := range strings.Split(string(output), "\n") {
			if line != "" {
				a.debug("%s", line)
			}
		}
	}

	if err != nil {
		return ExecError{Cmd: gitCmd, Args: args, Err: err}
	}
	return nil
}

// makedir creates the target dir and all missing directories on the path.
func (a adept) makedir(dir string) error {
	a.novice.makedir(dir)
	err := os.MkdirAll(dir, dirMode)
	if err != nil {
		return IOError{Cmd: "make dir", Arg: dir, Err: err}
	}

	return nil
}

// copy copies the file to the target.
// The target directory must exist.
func (a adept) copy(from, to string) error {
	a.novice.copy(from, to)

	to, dir, recursive, err := examine(from, to)
	if err != nil {
		return err
	}

	// if a dir is returned, create it
	if dir != "" {
		err = a.makedir(dir)
		if err != nil {
			return err
		}
	}

	if !recursive {
		// copy file
		return a.copyFile(from, to)
	}

	// copy directory
	return filepath.WalkDir(from, func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			// internal error
			return IOError{Cmd: "WalkDirFunc", Arg: path, Err: err}
		}

		if d.IsDir() {
			// do not copy directory entries
			return nil
		}

		subPath, _ := filepath.Rel(from, path)
		realTo := filepath.Join(to, subPath)
		err = os.MkdirAll(filepath.Dir(realTo), dirMode)
		if err != nil {
			return IOError{Cmd: "make dir", Arg: filepath.Dir(realTo), Err: err}
		}

		return a.copyFile(path, realTo)
	})
}

// copyFile copies the file to the target.
// The target directory must exist.
func (a adept) copyFile(from, to string) error {

	srcFile, err := os.Open(from)
	if err != nil {
		return IOError{Cmd: "open", Arg: from, Err: err}
	}
	defer srcFile.Close()

	targetFile, err := os.Create(to)
	if err != nil {
		return IOError{Cmd: "create", Arg: to, Err: err}
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, srcFile)
	if err != nil {
		return IOError{Cmd: "copy", Arg: from + " - " + to, Err: err}
	}

	return nil
}

// examine checks the source ("from") and target ("to") of a copy call.
// It returns:
//
//   - the new target for the copy
//   - a directory name if one should be created (might be empty)
//   - a bool to indicate recursive copy
//   - an error if one occurred
//
// The source is determined with the following logic.
//
//   - If from does not exist, it returns an error.
//   - If from is a file, then recursive will be false.
//   - If from is a directory, then recursive will be true.
//
// The target is determined with the following logic.
//
// If the target exists:
//
//   - If the check of "to" fails, the error is returned.
//   - If to is a file, it is returned
//   - If to is directory, it appends the last element from 'from' to 'to'
//
// If the target does not exist:
//
//   - If the source is a directory or if the source is a file
//     and to ends with a slash, it is used a a directory
//   - Else it is used as target file.
//
// If a directory to the target should be created, dir contains the name.
func examine(from, to string) (target, dir string, recursive bool, err error) {

	fromInfo, err := os.Stat(from)
	if err != nil {
		return "", "", false, IOError{Cmd: "stat", Arg: from, Err: err}
	}

	targetExists, targetIsDir, err := examineTarget(to)
	if err != nil {
		return "", "", false, err
	}
	target = filepath.Clean(to)

	if fromInfo.IsDir() {
		if targetExists {
			// target dir exists: append last directory from 'from'
			target = filepath.Join(target, filepath.Base(from))
		}
		return target, target, true, nil
	}

	if !targetExists {
		dir = target
		if !targetIsDir {
			dir = filepath.Dir(target)
		}
	}
	// target  exists
	if targetIsDir {
		// to is a directory
		return filepath.Join(target, filepath.Base(from)), dir, false, nil
	}

	return target, dir, false, nil
}

// examineTarget examines the to element. It returns if it does exist
// and if it is or should be a directory.
func examineTarget(to string) (exist, isdir bool, err error) {

	toInfo, err := os.Stat(to)
	if err != nil {
		if !os.IsNotExist(err) {
			return false, false, IOError{Cmd: "stat", Arg: to, Err: err}
		}
		// to does not exist
		if strings.HasSuffix(to, string(os.PathSeparator)) {
			return false, true, nil
		}
		return false, false, nil
	}

	// to does exist
	if toInfo.IsDir() {
		return true, true, nil
	}

	return true, false, nil
}
