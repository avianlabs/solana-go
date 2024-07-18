package rpc

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

type IsBlockhashValidOpts struct {
	// Commitment (optional) level to check for validity for.
	Commitment CommitmentType

	// MinContextSlot (optional) is the minimum slot that the request can be
	// evaulated at.
	MinContextSlot *uint64
}

// Returns whether a blockhash is still valid or not
//
// **NEW: This method is only available in solana-core v1.9 or newer. Please use
// `getFeeCalculatorForBlockhash` for solana-core v1.8**
func (cl *Client) IsBlockhashValid(
	ctx context.Context,
	// Blockhash to be queried. Required.
	blockHash solana.Hash,

	// Commitment requirement. Optional.
	commitment CommitmentType,
) (out *IsValidBlockhashResult, err error) {
	return cl.IsBlockhashValidWithOpts(
		ctx, blockHash, IsBlockhashValidOpts{
			Commitment: commitment,
		},
	)
}

// Returns whether a blockhash is still valid or not
//
// **NEW: This method is only available in solana-core v1.9 or newer. Please use
// `getFeeCalculatorForBlockhash` for solana-core v1.8**
func (cl *Client) IsBlockhashValidWithOpts(
	ctx context.Context,
	// Blockhash to be queried. Required.
	blockHash solana.Hash,
	opts IsBlockhashValidOpts,
) (out *IsValidBlockhashResult, err error) {
	params := []interface{}{blockHash}
	if opts.Commitment != "" {
		params = append(params, M{"commitment": string(opts.Commitment)})
	}
	if opts.MinContextSlot != nil {
		params = append(params, M{"minContextSlot": *opts.MinContextSlot})
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "isBlockhashValid", params)
	return
}
