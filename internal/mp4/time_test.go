package mp4

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetCreationTime(t *testing.T) {
	// fileHasMetadata is from core_test.go.
	date, err := GetCreationTime(fileHasMetadata)
	require.NoError(t, err)
	checkDate(t, date)
}

func checkDate(t *testing.T, date time.Time) {
	yr, mon, dy := date.Date()
	assert.Equal(t, 2022, yr)
	assert.Equal(t, time.Month(8), mon)
	assert.Equal(t, 22, dy)
	hr, min, sec := date.Clock()
	assert.Equal(t, 0, hr)
	assert.Equal(t, 15, min)
	assert.Equal(t, 13, sec)
}
