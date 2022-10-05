package exif

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCreationTime(t *testing.T) {
	// fileHasEXIF and checkDate() are from core_test.go.
	date, err := GetCreationTime(fileHasEXIF)
	require.NoError(t, err)
	checkDate(t, date)
}
