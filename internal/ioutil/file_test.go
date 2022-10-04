package ioutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fileCant   = "////snoofus.txt"
	fileNoSuch = "../../testdata/goober.txt"
	fileOther  = "../../testdata/other.txt"
	fileSource = "../../testdata/some.txt"
	fileTarget = "../../testdata/some.tgt"
)

func TestCopyFile(t *testing.T) {
	if err := os.Remove(fileTarget); err != nil {
		require.ErrorIs(t, err, os.ErrNotExist)
	}
	defer func() { assert.NoError(t, os.Remove(fileTarget)) }()
	assert.NoError(t, CopyFile(fileSource, fileTarget))
	assert.ErrorIs(t, CopyFile(fileSource, fileTarget), ErrIdentical)
	assert.ErrorIs(t, CopyFile(fileSource, fileOther), errDifferent)
	assert.ErrorIs(t, CopyFile(fileNoSuch, fileCant), os.ErrNotExist)
	assert.ErrorIs(t, CopyFile(fileSource, fileCant), os.ErrPermission)
}
