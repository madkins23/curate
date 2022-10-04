package exif

import (
	"testing"
	"time"

	EXIF "github.com/dsoprea/go-exif/v3"
	"github.com/madkins23/go-utils/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fileNoEXIF  = "../../testdata/some.txt"
	fileHasEXIF = "../../testdata/DSCF0013.JPG"
)

func Test_getIndex(t *testing.T) {
	_, err := GetIndex(fileNoEXIF)
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
	checkDate(t, date, 2022, time.Month(8), 23, 2, 43, 16)
}

func checkDate(t *testing.T, date time.Time, year int, month time.Month, day, hour, minute, second int) {
	yr, mon, dy := date.Date()
	assert.Equal(t, year, yr)
	assert.Equal(t, month, mon)
	assert.Equal(t, day, dy)
	hr, min, sec := date.Clock()
	assert.Equal(t, hour, hr)
	assert.Equal(t, minute, min)
	assert.Equal(t, second, sec)
}
