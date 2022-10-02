package exif

import (
	"fmt"
	"time"

	EXIF "github.com/dsoprea/go-exif/v3"
	EXIFcommon "github.com/dsoprea/go-exif/v3/common"
)

// GetCreationTime acquires the creation time of a source with EXIF properties.
func GetCreationTime(source string) (time.Time, error) {
	const creationTimeID = 0x132
	var ct time.Time
	if value, err := getValue(source, creationTimeID); err != nil {
		return ct, fmt.Errorf("get EXIF value: %w", err)
	} else if whenStr, ok := value.(string); !ok {
		return ct, fmt.Errorf("creation date not string")
	} else if ct, err = time.Parse("2006:01:02 15:04:05", whenStr); err != nil {
		return ct, fmt.Errorf("parse creation date: %w", err)
	}

	// TODO: Parsed time was local, how to get to UTC or do we care?

	return ct, nil
}

// getIndex acquires an index of EXIF properties.
func getIndex(source string) (EXIF.IfdIndex, error) {
	var index EXIF.IfdIndex
	if rawExif, err := EXIF.SearchFileAndExtractExif(source); err != nil {
		return index, fmt.Errorf("getting EXIF from file: %w", err)
	} else if im, err := EXIFcommon.NewIfdMappingWithStandard(); err != nil {
		return index, fmt.Errorf("getting EXIF mapping: %w", err)
	} else {
		ti := EXIF.NewTagIndex()
		if _, index, err = EXIF.Collect(im, ti, rawExif); err != nil {
			return index, fmt.Errorf("getting EXIF index: %w", err)
		} else {
			return index, nil
		}
	}
}

// getValue acquires the value of an EXIF property specified by tag ID.
// Tag names may be different for different devices but IDs seem to be fixed.
// IDs are only found within the IFD/Exif subtree of the index object returned by getIndex.
//
// ID resource: https://exiftool.org/TagNames/EXIF.html
func getValue(source string, tagID uint16) (interface{}, error) {
	index, err := getIndex(source)
	if err != nil {
		return nil, fmt.Errorf("get index: %w", err)
	}
	tagResults, err := index.RootIfd.FindTagWithId(tagID)
	if err != nil {
		tagResults, err = index.Lookup["IFD/Exif"].FindTagWithId(tagID)
	}
	if err != nil {
		return "", fmt.Errorf("lookup ID %d: %w", tagID, err)
	}
	if len(tagResults) != 1 {
		return "", fmt.Errorf("wrong number of results: %d", len(tagResults))
	} else if value, err := tagResults[0].Value(); err != nil {
		return "", fmt.Errorf("get tag value: %w", err)
	} else {
		return value, nil
	}
}
