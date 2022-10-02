package name

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/madkins23/curate/internal/exif"
	"github.com/madkins23/curate/internal/mp4"
)

const fmtGoogleTimestmap = "20060102_150405"

var errBadExtension = errors.New("unrecognized file extension")

// makeTimeStamp gets the file's creation date and returns a Google-format timestamp string.
func makeTimestamp(source string) (string, error) {
	//var timestamp string
	var creationTime time.Time
	var err error
	switch ext := strings.ToLower(filepath.Ext(source)); ext {
	case ".jpg", ".jpeg":
		if creationTime, err = exif.GetCreationTime(source); err != nil {
			return "", fmt.Errorf("get EXIF creation time: %w", err)
		}
	case ".mp4":
		if creationTime, err = mp4.GetCreationTime(source); err != nil {
			return "", fmt.Errorf("get MP4 creation time: %w", err)
		}
	default:
		return "", errBadExtension
	}

	return creationTime.Format(fmtGoogleTimestmap), nil
}
