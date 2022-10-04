package exif

import (
	"fmt"
	"time"
)

const (
	fmtCreationDate = "2006:01:02 15:04:05"
	idCreationDate  = 0x132
)

// GetCreationTime acquires the creation time of a source with EXIF properties.
func GetCreationTime(source string) (time.Time, error) {
	var ct time.Time

	index, err := GetIndex(source)
	if err != nil {
		return ct, fmt.Errorf("get EXIF index: %w", err)
	}

	var ok bool
	var value interface{}

	var dateTimeStr string
	if value, err = getValue(index, idCreationDate); err != nil {
		return ct, fmt.Errorf("get creation date value: %w", err)
	} else if dateTimeStr, ok = value.(string); !ok {
		return ct, fmt.Errorf("creation date not string")
	} else if ct, err = time.Parse(fmtCreationDate, dateTimeStr); err != nil {
		return ct, fmt.Errorf("parse creation date: %w", err)
	}

	// Don't look for sub-second value, it's probably not there and will be ignored anyway.
	// Don't look for time zone value, it's probably not there.

	return ct, nil
}
