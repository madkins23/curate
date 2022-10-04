package mp4

import (
	"fmt"
	"os"
	"time"

	"github.com/abema/go-mp4"
)

// GetCreationTime acquires the creation time of an MP4 source.
func GetCreationTime(source string) (time.Time, error) {
	var ct time.Time

	if metadata, err := getMetadata(source); err != nil {
		return ct, fmt.Errorf("get metadata: %w", err)
	} else if len(metadata) != 1 {
		return ct, fmt.Errorf("wrong number of metadata results (%d != 1): %w", len(metadata), err)
	} else if payload, ok := metadata[0].Payload.(*mp4.Mvhd); !ok {
		return ct, fmt.Errorf("convert metadata payload to mvhd: %w", err)
	} else {
		// Property mvhd/CreationTimeV0 is seconds since Jan 1, 1904 (UTC) for some reason.
		// Don't look for sub-second value, it's probably not there and will be ignored anyway.
		ct = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).
			Add(time.Second * time.Duration(payload.CreationTimeV0)).
			// Convert to local time zone.
			In(time.Local)
	}

	return ct, nil
}

// getMetadata returns the moov/mvhd section of the MP4 file.
// This section contains various properties of the video.
//
// Spec resource: https://www.cimarronsystems.com/wp-content/uploads/2017/04/Elements-of-the-H.264-VideoAAC-Audio-MP4-Movie-v2_0.pdf
func getMetadata(path string) ([]*mp4.BoxInfoWithPayload, error) {
	if file, err := os.Open(path); err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	} else {
		return mp4.ExtractBoxWithPayload(file, nil,
			mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypeMvhd()})
	}
}
