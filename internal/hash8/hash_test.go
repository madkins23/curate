package hash8

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Taken from https://crccalc.com/:
	testData       = "123456789"
	testHash uint8 = 0xA1
)

func Test_Hash8(t *testing.T) {
	hash8 := New()
	require.NotNil(t, hash8)
	size, err := hash8.Write([]byte(testData))
	require.NoError(t, err)
	assert.Equal(t, len(testData), size)
	assert.Equal(t, testHash, hash8.Sum8())
}
