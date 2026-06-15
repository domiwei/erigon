package solid

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPendingConsolidationJSONUint64AsString(t *testing.T) {
	c := &PendingConsolidation{SourceIndex: 1, TargetIndex: 2}
	got, err := json.Marshal(c)
	require.NoError(t, err)
	require.Contains(t, string(got), `"source_index":"1"`)
	require.Contains(t, string(got), `"target_index":"2"`)
}
