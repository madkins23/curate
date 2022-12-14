package name

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	crc8hash "github.com/madkins23/go-hash/pkg/sigurn/crc8"
	"github.com/sigurn/crc8"

	"github.com/madkins23/curate/internal/exif"
	"github.com/madkins23/curate/internal/mp4"
)

const fmtGoogleTimestmap = "20060102_150405"

var errBadExtension = errors.New("unrecognized file extension")

// makeTimeStamp gets the file's creation date and returns a Google-format timestamp string.
// Timestamps include random numbers for milliseconds to avoid potential collisions.
// While some media files may include sub-second values, many do not.
func makeTimestamp(source string) (string, error) {
	//var timestamp string
	var creationTime time.Time
	var err error
	switch ext := strings.ToLower(filepath.Ext(source)); ext {
	case ".jpg", ".jpeg":
		if creationTime, err = exif.GetCreationTime(source); err != nil {
			return "", fmt.Errorf("get EXIF creation time: %w", err)
		}
	case ".mov", ".mp4":
		if creationTime, err = mp4.GetCreationTime(source); err != nil {
			return "", fmt.Errorf("get MP4 creation time: %w", err)
		}
	default:
		return "", errBadExtension
	}

	var millis string
	if millis, err = makeCRC8(source); err != nil {
		// Use CRC8 to generate a three digit fake millisecond value.
		// Using the CRC will result in a constant number for the same file
		// in case it is processed more than once for some reason.
		return "", fmt.Errorf("get CRC8 for millis: %w", err)
	}

	return creationTime.Format(fmtGoogleTimestmap) + millis, nil
}

// Make a CRC8 string of three digits for the contents of the specified source file.
// The resulting number will be 000..256, formatted as three digit decimal string with leading zeros.
func makeCRC8(source string) (string, error) {
	if file, err := os.Open(source); err != nil {
		return "", fmt.Errorf("open file: %w", err)
	} else {
		defer func() { _ = file.Close() }()
		hash := crc8hash.New(crc8.CRC8_MAXIM)
		if _, err = io.Copy(hash, file); err != nil {
			return "", fmt.Errorf("calculate hash: %w", err)
		}
		return fmt.Sprintf("%03d", hash.Sum8()), nil
	}
}
