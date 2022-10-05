package mp4

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fileCant        = "////snoofus.txt"
	fileNoSuch      = "../../testdata/goober.txt"
	fileNoMetadata  = "../../testdata/some.txt"
	fileHasMetadata = "../../testdata/DSCF0004.MP4"
)

func Test_getMetadata(t *testing.T) {
	_, err := getMetadata(fileCant)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = getMetadata(fileNoSuch)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = getMetadata(fileNoMetadata)
	assert.ErrorIs(t, err, errNoMetadata)
	meta, err := getMetadata(fileHasMetadata)
	require.NoError(t, err)
	require.NotNil(t, meta)
	// Test CreationTimeV0 because that's the field we currently need.
	assert.NotZero(t, meta.CreationTimeV0)
}
