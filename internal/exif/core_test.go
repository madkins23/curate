package exif

import (
	"os"
	"testing"
	"time"

	EXIF "github.com/dsoprea/go-exif/v3"
	"github.com/madkins23/go-utils/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fileCant    = "////snoofus.txt"
	fileNoSuch  = "../../testdata/goober.txt"
	fileNoEXIF  = "../../testdata/some.txt"
	fileHasEXIF = "../../testdata/DSCF0013.JPG"
)

func Test_getIndex(t *testing.T) {
	_, err := GetIndex(fileCant)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = GetIndex(fileNoSuch)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = GetIndex(fileNoEXIF)
	assert.ErrorIs(t, err, EXIF.ErrNoExif)
	index, err := GetIndex(fileHasEXIF)
	assert.NoError(t, err)
	assert.NotNil(t, index)
	assert.NotNil(t, index.RootIfd)
	log.Console()
	assert.NoError(t, EnumerateIndex(index))
}

func Test_getValue(t *testing.T) {
	index, err := GetIndex(fileHasEXIF)
	require.NoError(t, err)
	require.NotNil(t, index)
	// Test Creation Date because that's the one we currently need.
	value, err := getValue(index, idCreationDate)
	require.NoError(t, err)
	valStr, ok := value.(string)
	require.True(t, ok)
	date, err := time.Parse(fmtCreationDate, valStr)
	require.NoError(t, err)
	checkDate(t, date)
}

func checkDate(t *testing.T, date time.Time) {
	yr, mon, dy := date.Date()
	assert.Equal(t, 2022, yr)
	assert.Equal(t, time.Month(8), mon)
	assert.Equal(t, 23, dy)
	hr, min, sec := date.Clock()
	assert.Equal(t, 2, hr)
	assert.Equal(t, 43, min)
	assert.Equal(t, 16, sec)
}
