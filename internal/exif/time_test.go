package exif

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const filePXL = "../../testdata/PXL_20220215_002116112.jpg"

func TestGetCreationTime(t *testing.T) {
	date, err := GetCreationTime(filePXL)
	require.NoError(t, err)
	// From core_test.go:
	checkDate(t, date, 2022, time.Month(2), 14, 16, 21, 16)
}
