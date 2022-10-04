package exif

import (
	"fmt"
	"strconv"

	EXIF "github.com/dsoprea/go-exif/v3"
	EXIFcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/rs/zerolog/log"
)

func EnumerateIndex(index EXIF.IfdIndex) error {
	err := index.RootIfd.EnumerateTagsRecursively(func(ifd *EXIF.Ifd, ite *EXIF.IfdTagEntry) error {
		event :=
			log.Debug().
				Str("path", ite.IfdPath()+"/"+ite.TagName()).
				Str("ID", "0x"+strconv.FormatUint(uint64(ite.TagId()), 16))
		value, err := ite.Value()
		if err == nil {
			event.Interface("value", value).Msg("Item")
		} else {
			event.Err(err).Msg("Error")
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// GetIndex acquires an index of EXIF properties.
func GetIndex(source string) (EXIF.IfdIndex, error) {
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
// IDs are only found within the IFD/Exif subtree of the index object returned by GetIndex.
//
// ID resource: https://exiftool.org/TagNames/EXIF.html
func getValue(index EXIF.IfdIndex, tagID uint16) (interface{}, error) {
	tagResults, err := index.RootIfd.FindTagWithId(tagID)
	if err != nil {
		tagResults, err = index.Lookup["IFD/Exif"].FindTagWithId(tagID)
	}
	if err != nil {
		//_ = EnumerateIndex(index)
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
