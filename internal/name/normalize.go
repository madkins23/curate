package name

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
)

var (
	ErrNoPattern = errors.New("basename matches no known pattern")
	errNoPeriods = errors.New("basename has no periods")
	errTBD       = errors.New("basename pattern not yet implemented")
)

// Normalize checks a file's basename of the source path for format and uniqueness,
// returning a better name if necessary, otherwise just the basename.
func Normalize(source string) (string, error) {
	basename := path.Base(source)
	chunks := strings.Split(basename, ".")
	if len(chunks) < 2 {
		return "", errNoPeriods
	}
	name := []byte(chunks[0])

	if err := initPatterns(); err != nil {
		return "", fmt.Errorf("initialize initPatterns: %w", err)
	}

	if ptnGoogle.Match(name) {
		// Just return the basename
	} else if ptnDSC.Match(name) {
		if timestamp, err := makeTimestamp(source); err != nil {
			return "", fmt.Errorf("making timestamp: %w", err)
		} else {
			chunks[0] = "DSC_" + timestamp
			return strings.Join(chunks, "."), nil
		}
	} else {
		return basename, ErrNoPattern
	}

	return basename, nil
}

var (
	ptnInitialized bool
	ptnDSC         *regexp.Regexp
	ptnGoogle      *regexp.Regexp
)

// initPatterns initializes the recognized patterns for file names.
func initPatterns() error {
	if !ptnInitialized {
		var err error
		if ptnGoogle, err = regexp.Compile("[[:alpha:]]+_\\d{8}_\\d{9}"); err != nil {
			return fmt.Errorf("compile Google pattern: %w", err)
		} else if ptnDSC, err = regexp.Compile("DSCF\\d+"); err != nil {
			return fmt.Errorf("compile DSC pattern: %w", err)
		}
		ptnInitialized = true
	}
	return nil
}
