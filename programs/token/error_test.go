package token

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

func TestParseError(t *testing.T) {
	t.Parallel()
	var (
		err error
		ok  bool
	)

	tx, err := solana.NewTransaction([]solana.Instruction{
		NewTransferInstruction(
			0,
			solana.PublicKey{},
			solana.PublicKey{},
			solana.PublicKey{},
			nil,
		).Build(),
	}, solana.Hash{})
	require.NoError(t, err)
	err, ok = solana.ParseTransactionError(tx, map[string]interface{}{
		"InstructionError": []interface{}{
			float64(0),
			map[string]interface{}{
				"Custom": float64(1),
			},
		},
	})
	require.True(t, ok)
	assert.Equal(t, &solana.TransactionError_InstructionError{
		Index: 0,
		Cause: &solana.InstructionError_Custom{
			Code:  1,
			Cause: Error_InsufficientFunds{},
		},
	}, err)
}
