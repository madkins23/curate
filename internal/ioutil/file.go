package ioutil

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/udhos/equalfile"
)

var (
	ErrIdentical = errors.New("identical files")
	fileCompare  = equalfile.New(nil, equalfile.Options{})
)

func CopyFile(source, target string) error {
	if _, err := os.Stat(target); err == nil {
		if equal, err := fileCompare.CompareFile(source, target); err != nil {
			return fmt.Errorf("compare files: %w", err)
		} else if equal {
			return ErrIdentical
		} else {
			return fmt.Errorf("pre-existing file not identical")
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
