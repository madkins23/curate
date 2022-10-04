package ioutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dirExists      = "../../testdata"
	dirNonExistent = "no-such-directory"
	dirNotDir      = "../../internal/ioutil/dir.go"
)

func TestCheckDir(t *testing.T) {
	assert.NoError(t, CheckDir(dirExists))
	assert.ErrorIs(t, CheckDir(dirNonExistent), errNonExistent)
	assert.ErrorIs(t, CheckDir(dirNotDir), errNotDirectory)
}
