package cltypes_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/erigontech/erigon/cl/clparams"
	"github.com/erigontech/erigon/cl/cltypes"
	"github.com/erigontech/erigon/cl/cltypes/solid"
	"github.com/erigontech/erigon/common"
)

func TestExecutionRequestsEncodingSizeSSZ(t *testing.T) {
	cfg := &clparams.MainnetBeaconConfig

	t.Run("electra_empty", func(t *testing.T) {
		er := cltypes.NewExecutionRequestsWithVersion(cfg, clparams.ElectraVersion)
		encoded, err := er.EncodeSSZ(nil)
		require.NoError(t, err)
		require.Equal(t, len(encoded), er.EncodingSizeSSZ())
	})

	t.Run("electra_populated", func(t *testing.T) {
		er := cltypes.NewExecutionRequestsWithVersion(cfg, clparams.ElectraVersion)
		er.Deposits.Append(&solid.DepositRequest{
			PubKey: common.Bytes48{1},
			Amount: 32_000_000_000,
		})
		er.Withdrawals.Append(&solid.WithdrawalRequest{
			SourceAddress: common.Address{2},
			Amount:        1_000_000_000,
		})
		er.Consolidations.Append(&solid.ConsolidationRequest{
			SourceAddress: common.Address{3},
		})
		encoded, err := er.EncodeSSZ(nil)
		require.NoError(t, err)
		require.Equal(t, len(encoded), er.EncodingSizeSSZ())
	})

	t.Run("gloas_empty", func(t *testing.T) {
		er := cltypes.NewExecutionRequestsWithVersion(cfg, clparams.GloasVersion)
		encoded, err := er.EncodeSSZ(nil)
		require.NoError(t, err)
		require.Equal(t, len(encoded), er.EncodingSizeSSZ())
	})

	t.Run("gloas_populated", func(t *testing.T) {
		er := cltypes.NewExecutionRequestsWithVersion(cfg, clparams.GloasVersion)
		er.Deposits.Append(&solid.DepositRequest{
			PubKey: common.Bytes48{1},
			Amount: 32_000_000_000,
		})
		er.Withdrawals.Append(&solid.WithdrawalRequest{
			SourceAddress: common.Address{2},
			Amount:        1_000_000_000,
		})
		er.Consolidations.Append(&solid.ConsolidationRequest{
			SourceAddress: common.Address{3},
		})
		er.BuilderDeposits.Append(&solid.BuilderDepositRequest{
			PubKey: common.Bytes48{4},
			Amount: 100_000_000,
		})
		er.BuilderExits.Append(&solid.BuilderExitRequest{
			SourceAddress: common.Address{5},
		})
		encoded, err := er.EncodeSSZ(nil)
		require.NoError(t, err)
		require.Equal(t, len(encoded), er.EncodingSizeSSZ())
	})
}
