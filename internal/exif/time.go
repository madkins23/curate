package exif

import (
	"fmt"
	"time"

	EXIF "github.com/dsoprea/go-exif/v3"
	EXIFcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/rs/zerolog/log"
)

const (
	fmtCreationDate = "2006:01:02 15:04:05"
	idCreationDate  = 0x132
	idTimeZone      = 0x882a
)

// GetCreationTime acquires the creation time of a source with EXIF properties.
func GetCreationTime(source string) (time.Time, error) {
	var ct time.Time

	index, err := getIndex(source)
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

	if value, err = getValue(index, idTimeZone); err != nil {
		// Time zone not always available, just skip it
	} else {
		// Just in case one actually shows up:
		log.Debug().Interface("time zone(s)", value).Msg("Time zone available!")
	}

	return ct, nil
}

//func enumerateIndex(index EXIF.IfdIndex) error {
//	err := index.RootIfd.EnumerateTagsRecursively(func(ifd *EXIF.Ifd, ite *EXIF.IfdTagEntry) error {
//		log.Debug().Str("path", ite.IfdPath()+"/"+ite.TagName()).
//			Str("ID", "0x"+strconv.FormatUint(uint64(ite.TagId()), 16)).Msg("tag")
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

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
func getValue(index EXIF.IfdIndex, tagID uint16) (interface{}, error) {
	tagResults, err := index.RootIfd.FindTagWithId(tagID)
	if err != nil {
		tagResults, err = index.Lookup["IFD/Exif"].FindTagWithId(tagID)
	}
	if err != nil {
		//_ = enumerateIndex(index)
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
