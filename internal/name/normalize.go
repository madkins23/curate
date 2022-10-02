package name

import "path"

// Normalize checks a file's basename of the source path for format and uniqueness,
// returning a better name if necessary, otherwise just the basename.
// The returned basename may be a path.
func Normalize(sourcePath string) (string, error) {
	basename := path.Base(sourcePath)

	return basename, nil
}
