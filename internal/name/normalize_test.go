package name

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Some constants are defined in timestamp_test.go.

const (
	filePXLnotExist = "../../testdata/PXL_20220215_009916662.jpg"
	fileJPGname     = "DSC_" + fileJPGts + ".JPG"
	fileMP4name     = "DSC_" + fileMP4ts + ".MP4"
)

func TestNormalize(t *testing.T) {
	// Must be first thing before Normalize() is ever executed.
	assert.True(t, testInitPatterns(t))

	_, err := Normalize(fileBadExt)
	assert.ErrorIs(t, err, ErrNoPattern)
	_, err = Normalize(fileCant)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = Normalize(fileNoSuch)
	assert.ErrorIs(t, err, os.ErrNotExist)
	_, err = Normalize(filePXLnotExist)
	assert.ErrorIs(t, err, os.ErrNotExist)
	name, err := Normalize(fileJPG)
	require.NoError(t, err)
	assert.Equal(t, fileJPGname, name)
	name, err = Normalize(fileMP4)
	require.NoError(t, err)
	assert.Equal(t, fileMP4name, name)
}

func testInitPatterns(t *testing.T) bool {
	// Can only run this once per test run.
	if !rgxInitialized {
		assert.Nil(t, rgxDSC)
		assert.Nil(t, rgxGoogle)
		require.NoError(t, initPatterns())
		require.True(t, rgxInitialized)
		require.NotNil(t, rgxDSC)
		require.NotNil(t, rgxGoogle)
		return true
	}
	return false
}
