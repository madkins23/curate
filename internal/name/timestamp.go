package name

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/madkins23/go-utils/msg"

	"github.com/madkins23/curate/internal/exif"
)

const fmtGoogleTimestmap = "20060102_150405"

var errBadExtension = errors.New("unrecognized file extension")

// makeTimeStamp gets the file's creation date and returns a Google-format timestamp string.
func makeTimestamp(source string) (string, error) {
	//var timestamp string
	switch ext := strings.ToLower(filepath.Ext(source)); ext {
	case ".jpg", ".jpeg":
		if creationTime, err := exif.GetCreationTime(source); err != nil {
			return "", fmt.Errorf("get creation time: %w", err)
		} else {
			return creationTime.Format(fmtGoogleTimestmap), nil
		}
	case ".mp4":
		return "", &msg.ErrNotImplemented{Name: "MP4"}
	default:
		return "", errBadExtension
	}

	//return timestamp, nil
}
