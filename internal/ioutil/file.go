package ioutil

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/udhos/equalfile"
)

var (
	// ErrIdentical is returned if a pre-existing file is identical to the new file.
	ErrIdentical = errors.New("identical files")

	// errDifferent is returned if a pre-existing file is different from the new file.
	errDifferent = errors.New("pre-existing file not identical")

	// fileCompare object compares two files using github.com/udhos/equalfile.
	fileCompare = equalfile.New(nil, equalfile.Options{})
)

// CopyFile copies the source file to the specified target path,
// checking for pre-existing files with the target path.
// ErrIdentical is returned if a pre-existing file is identical to the new file.
func CopyFile(source, target string) error {
	if _, err := os.Stat(target); err == nil {
		if equal, err := fileCompare.CompareFile(source, target); err != nil {
			return fmt.Errorf("compare files: %w", err)
		} else if equal {
			return ErrIdentical
		} else {
			return errDifferent
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err := copyFile(source, target); err != nil {
			return fmt.Errorf("copy file: %w", err)
		}
	} else {
		return fmt.Errorf("stat target file: %w", err)
	}
	return nil
}

// copyFile copies the source file to the specified target path.
// No checks are done for pre-existing files with the target path.
func copyFile(source, target string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer func() { _ = sourceFile.Close() }()
	targetFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("create target file: %w", err)
	}
	defer func() { _ = targetFile.Close() }()
	if _, err = io.Copy(targetFile, sourceFile); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}
	return nil
}
