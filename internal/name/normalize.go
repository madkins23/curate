package name

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	ErrNoPattern = errors.New("basename matches no known pattern")
	errNoPeriods = errors.New("basename has no periods")
	rgxDSC       *regexp.Regexp
	rgxGoogle    *regexp.Regexp
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

	if rgxGoogle.Match(name) {
		// Do this check so that results will be consistent.
		// Other branch(es?) will attempt to access the source,
		// resulting in an error if the source does not exist.
		if _, err := os.Stat(source); err != nil {
			return "", fmt.Errorf("check existence: %w", err)
		}
		// Just return the basename
	} else if rgxDSC.Match(name) {
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

const (
	ptnDSC    = "^DSC[[:alpha:]]\\d+$"
	ptnGoogle = "^[[:alpha:]]+_\\d{8}_\\d{9}$"
)

var rgxInitialized bool

// initPatterns initializes the recognized patterns for file names.
func initPatterns() error {
	if !rgxInitialized {
		var err error
		if rgxGoogle, err = regexp.Compile(ptnGoogle); err != nil {
			return fmt.Errorf("compile Google pattern: %w", err)
		} else if rgxDSC, err = regexp.Compile(ptnDSC); err != nil {
			return fmt.Errorf("compile DSC pattern: %w", err)
		}
		rgxInitialized = true
	}
	return nil
}
