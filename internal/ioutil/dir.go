package ioutil

import (
	"errors"
	"fmt"
	"os"
)

var (
	errNotDirectory = errors.New("not a directory")
	errNonExistent  = errors.New("directory does not exist")
)

// CheckDir checks to see if the specified directory exists and if not returns an error.
func CheckDir(dir string) error {
	if stat, err := os.Stat(dir); err == nil {
		if !stat.IsDir() {
			return errNotDirectory
		}
	} else if errors.Is(err, os.ErrNotExist) {
		return errNonExistent
	} else {
		return fmt.Errorf("stat dir error: %w", err)
	}
	return nil
}
