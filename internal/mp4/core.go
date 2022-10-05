package mp4

import (
	"errors"
	"fmt"
	"os"

	"github.com/abema/go-mp4"
)

var (
	errbadPayload = errors.New("bad payload")
	errNoMetadata = errors.New("no metadata")
)

// getMetadata returns the moov/mvhd section of the MP4 file.
// This section contains various properties of the video.
//
// Spec resource: https://www.cimarronsystems.com/wp-content/uploads/2017/04/Elements-of-the-H.264-VideoAAC-Audio-MP4-Movie-v2_0.pdf
func getMetadata(path string) (*mp4.Mvhd, error) {
	if file, err := os.Open(path); err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	} else {
		defer func() { _ = file.Close() }()
		box, err := mp4.ExtractBoxWithPayload(file, nil,
			mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypeMvhd()})
		if err != nil {
			return nil, fmt.Errorf("extract moov/mvhd: %w", err)
		} else if len(box) < 1 {
			return nil, errNoMetadata
		} else if mvhd, ok := box[0].Payload.(*mp4.Mvhd); !ok {
			return nil, errbadPayload
		} else {
			return mvhd, nil
		}
	}
}
