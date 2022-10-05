package name

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fileBadExt  = "../../testdata/some.txt"
	fileCant    = "////DSCF0666.jpg"
	fileNoSuch  = "../../testdata/DSCF0666.MP4"
	fileJPG     = "../../testdata/DSCF0013.JPG"
	fileJPGhash = "023"
	fileJPGts   = "20220823_024316023"
	fileMP4     = "../../testdata/DSCF0004.MP4"
	fileMP4hash = "061"
	fileMP4ts   = "20220822_001513061"
)

func Test_makeTimestamp(t *testing.T) {
	_, err := makeTimestamp(fileBadExt)
	assert.ErrorIs(t, err, errBadExtension)
	_, err = makeTimestamp(fileCant)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = makeTimestamp(fileNoSuch)
	assert.ErrorIs(t, err, os.ErrNotExist)
	name, err := makeTimestamp(fileJPG)
	require.NoError(t, err)
	assert.Equal(t, fileJPGts, name)
	name, err = makeTimestamp(fileMP4)
	require.NoError(t, err)
	assert.Equal(t, fileMP4ts, name)
}

func Test_makeCRC8(t *testing.T) {
	_, err := makeCRC8(fileCant)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = makeCRC8(fileNoSuch)
	assert.ErrorIs(t, err, os.ErrNotExist)
	crc, err := makeCRC8(fileJPG)
	require.NoError(t, err)
	assert.Equal(t, fileJPGhash, crc)
	crc, err = makeCRC8(fileMP4)
	require.NoError(t, err)
	assert.Equal(t, fileMP4hash, crc)
}
