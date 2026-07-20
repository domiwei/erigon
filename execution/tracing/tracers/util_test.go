package tracers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMemoryCopyPaddedOverflowStopsBeforeAllocation(t *testing.T) {
	// Offsets derived from the glamsterdam-devnet-6 CREATE2 that triggered
	// "makeslice: len out of range" via signed int64 overflow.
	_, err := GetMemoryCopyPadded(nil, 8454290623495228678, 1395224211451850339)

	require.Error(t, err)
	require.Contains(t, err.Error(), "overflows int64")
}

func TestGetMemoryCopyPaddedNormalOperation(t *testing.T) {
	m := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	// Fully inside memory
	result, err := GetMemoryCopyPadded(m, 1, 3)
	require.NoError(t, err)
	require.Equal(t, []byte{0x02, 0x03, 0x04}, result)

	// Extends beyond memory (zero-padded)
	result, err = GetMemoryCopyPadded(m, 3, 4)
	require.NoError(t, err)
	require.Equal(t, []byte{0x04, 0x05, 0x00, 0x00}, result)

	// Negative offset
	_, err = GetMemoryCopyPadded(m, -1, 3)
	require.Error(t, err)
}
