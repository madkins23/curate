package mp4

import (
	"fmt"
	"os"
	"time"

	"github.com/abema/go-mp4"
)

func GetCreationTime(source string) (time.Time, error) {
	var ct time.Time

	if metadata, err := getMetadata(source); err != nil {
		return ct, fmt.Errorf("get metadata: %w")
	} else if len(metadata) != 1 {
		return ct, fmt.Errorf("wrong number of metadata results (%d != 1): %w", len(metadata), err)
	} else if payload, ok := metadata[0].Payload.(*mp4.Mvhd); !ok {
		return ct, fmt.Errorf("convert metadata payload to mvhd: %w", err)
	} else {
		// Mvhd/CreationTimeV0 is seconds since Jan 1, 1904 for some reason.
		ct = time.Date(1904, time.January, 1, 0, 0, 0, 0, time.UTC).
			Add(time.Second * time.Duration(payload.CreationTimeV0))
	}

	return ct, nil
}

func getMetadata(path string) ([]*mp4.BoxInfoWithPayload, error) {
	if file, err := os.Open(path); err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	} else {
		return mp4.ExtractBoxWithPayload(file, nil,
			mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypeMvhd()})
	}
}
