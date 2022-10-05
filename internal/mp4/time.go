package mp4

import (
	"fmt"
	"time"
)

// GetCreationTime acquires the creation time of an MP4 source.
func GetCreationTime(source string) (time.Time, error) {
	var ct time.Time

	if metadata, err := getMetadata(source); err != nil {
		return ct, fmt.Errorf("get metadata: %w", err)
	} else {
		// Property mvhd/CreationTimeV0 is seconds since Jan 1, 1904 (UTC) for some reason.
		// Don't look for sub-second value, it's probably not there and will be ignored anyway.
		ct = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).
			Add(time.Second * time.Duration(metadata.CreationTimeV0)).
			// Convert to local time zone.
			In(time.Local)
	}

	return ct, nil
}
