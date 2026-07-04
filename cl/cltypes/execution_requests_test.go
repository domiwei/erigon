package cltypes

import (
	"testing"

	"github.com/erigontech/erigon/cl/clparams"
	"github.com/stretchr/testify/require"
)

func TestExecutionRequests_EncodingSizeSSZ(t *testing.T) {
	cfg := clparams.MainnetBeaconConfig
	er := NewExecutionRequests(&cfg)

	encoded, err := er.EncodeSSZ(nil)
	require.NoError(t, err)

	require.Equal(t, len(encoded), er.EncodingSizeSSZ(),
		"EncodingSizeSSZ() must match actual encoded length")
}
