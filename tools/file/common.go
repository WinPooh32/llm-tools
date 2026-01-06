package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const escapesFromParentErr = "can only access files and directories beneath the current working directory"

var errEscapesFromParent = errors.New("escapes from parent directory")

func openFile(baseDir string, name string, ro bool) (*os.File, error) {
	var (
		rel string
		err error
	)

	if filepath.IsLocal(name) {
		rel = name
	} else {
		rel, err = filepath.Rel(baseDir, name)
		if err != nil {
			return nil, errEscapesFromParent
		}
	}

	r, err := os.OpenRoot(baseDir)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var flag int
	if ro {
		flag = os.O_RDONLY
	} else {
		flag = os.O_RDWR
	}

	file, err := r.OpenFile(rel, flag, 0666)
	if err != nil {
		return nil, fmt.Errorf("open file: os.Root: %w", err)
	}

	return file, nil
}
