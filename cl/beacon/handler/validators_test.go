package handler

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorIdentityResponseJSONUint64AsString(t *testing.T) {
	v := &validatorIdentityResponse{Index: 42, ActivationEpoch: 100}
	got, err := json.Marshal(v)
	require.NoError(t, err)
	require.Contains(t, string(got), `"index":"42"`)
	require.Contains(t, string(got), `"activation_epoch":"100"`)
}
