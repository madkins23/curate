package exif

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	EXIF "github.com/dsoprea/go-exif/v3"
	EXIFcommon "github.com/dsoprea/go-exif/v3/common"
)

// GetCreationTime acquires the creation time of a source with EXIF properties.
func GetCreationTime(source string) (time.Time, error) {
	var ct time.Time

	index, err := getIndex(source)
	if err != nil {
		return ct, fmt.Errorf("get EXIF index: %w", err)
	}

	const IDcreationDate = 0x132
	const IDsubSeconds = 0x9292
	const IDtimeZone = 0x882a

	var ok bool
	var value interface{}

	var dateTimeStr string
	if value, err = getValue(index, IDcreationDate); err != nil {
		return ct, fmt.Errorf("get creation date value: %w", err)
	} else if dateTimeStr, ok = value.(string); !ok {
		return ct, fmt.Errorf("creation date not string")
	}

	var subSecsStr string
	if value, err = getValue(index, IDsubSeconds); err != nil {
		// Subseconds not always available, make up random millis
		subSecsStr = fmt.Sprintf("%03d", rand.Intn(1000))
	} else if subSecsStr, ok = value.(string); !ok {
		return ct, fmt.Errorf("subseconds not string")
	} else if ptnJustDigits, err := regexp.Compile("\\d{3}"); err != nil {
		return ct, fmt.Errorf("compile digits pattern: %w", err)
	} else if !ptnJustDigits.Match([]byte(subSecsStr)) {
		return ct, fmt.Errorf("subseconds '%s' don't match pattern", subSecsStr)
	}

	if ct, err = time.Parse("2006:01:02 15:04:05.000", dateTimeStr+"."+subSecsStr); err != nil {
		return ct, fmt.Errorf("parse creation date: %w", err)
	}

	// TODO: Parsed time was local, how to get to UTC or do we care?
	if value, err = getValue(index, IDtimeZone); err != nil {
		// Time zone not always available, just skip it
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
