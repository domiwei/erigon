package solid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDepositRequestJSONUint64AsString(t *testing.T) {
	d := &DepositRequest{Amount: 32000000000, Index: 2457249}
	got, err := json.Marshal(d)
	require.NoError(t, err)
	require.Contains(t, string(got), `"amount":"32000000000"`)
	require.Contains(t, string(got), `"index":"2457249"`)
}

func TestPendingDepositJSONUint64AsString(t *testing.T) {
	d := &PendingDeposit{Amount: 32000000000, Slot: 13605852}
	got, err := json.Marshal(d)
	require.NoError(t, err)
	require.Contains(t, string(got), `"amount":"32000000000"`)
	require.Contains(t, string(got), `"slot":"13605852"`)
}
