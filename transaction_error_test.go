package solana

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTransactionError(t *testing.T) {
	t.Parallel()
	var (
		err error
		ok  bool
	)

	err, ok = ParseTransactionError(nil, "AccountInUse")
	require.True(t, ok)
	require.ErrorIs(t, err, TransactionError_AccountInUse{})

	err, ok = ParseTransactionError(nil, map[string]interface{}{
		"InstructionError": []interface{}{
			float64(1),
			"InvalidArgument",
		},
	})
	require.True(t, ok)
	var e *TransactionError
	require.ErrorAs(t, err, &e)
	assert.Equal(t, &TransactionError_InstructionError{
		Index: 1,
		Cause: InstructionError_InvalidArgument{},
	}, err.(interface{ Unwrap() error }).Unwrap())

	err, ok = ParseTransactionError(nil, map[string]interface{}{
		"InstructionError": []interface{}{
			float64(1),
			map[string]interface{}{
				"Custom": float64(16),
			},
		},
	})
	assert.Equal(t, &TransactionError_InstructionError{
		Index: 1,
		Cause: &InstructionError_Custom{
			Code:  16,
			Cause: nil,
		},
	}, err.(interface{ Unwrap() error }).Unwrap())
}
