package alchemist

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// ListPages lists the formula file for the provided tasks.
// The tasklist is a task (directory) name or a comma separated list of task names.
// Each directory must contain a gitalchemy.yaml file.
//
// It returns an error if the info of the gitalchemy.yaml file can not be accessed.
func ListPages(cfgdir string, tasklist ...string) ([]string, error) {

	result := []string{}
	for _, task := range tasklist {
		fileName := filepath.Join(cfgdir, task, FormulaFileName)
		_, err := os.Stat(fileName)
		if err != nil {
			return result, IOError{Cmd: "stat", Arg: fileName, Err: err}
		}
		result = append(result, fileName)
	}

	return result, nil
}

// ListBookContent lists the gitalchemy.yaml files of all
// subdirectories of the cfgdir directory.
//
// It returns an error if it can not access the directories.
//
// It does not return an error if a directory does not contain
// a gitalchemy.yaml file, but this directory is not included
// in the returned list.
func ListBookContent(cfgdir string) ([]string, error) {

	var result []string
	entries, err := os.ReadDir(cfgdir)
	if err != nil {
		return result, IOError{Cmd: "read dir", Arg: cfgdir, Err: err}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			pages, err := ListPages(cfgdir, entry.Name())
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return result, fmt.Errorf("list book content: %w", err)
			}
			result = append(result, pages...)
		}
	}

	return result, nil
}
