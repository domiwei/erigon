package solid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithdrawalRequestJSONUint64AsString(t *testing.T) {
	w := &WithdrawalRequest{Amount: 798102025}
	got, err := json.Marshal(w)
	require.NoError(t, err)
	require.Contains(t, string(got), `"amount":"798102025"`)
}

func TestPendingPartialWithdrawalJSONUint64AsString(t *testing.T) {
	w := &PendingPartialWithdrawal{
		ValidatorIndex:    42,
		Amount:            1000000000,
		WithdrawableEpoch: 123456,
	}
	got, err := json.Marshal(w)
	require.NoError(t, err)
	require.Contains(t, string(got), `"validator_index":"42"`)
	require.Contains(t, string(got), `"amount":"1000000000"`)
	require.Contains(t, string(got), `"withdrawable_epoch":"123456"`)
}
