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
	var opts *IsBlockhashValidOpts
	if commitment != "" {
		opts = &IsBlockhashValidOpts{
			Commitment: commitment,
		}
	}
	return cl.IsBlockhashValidWithOpts(
		ctx, blockHash, opts,
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
	opts *IsBlockhashValidOpts,
) (out *IsValidBlockhashResult, err error) {
	obj := M{}
	if opts != nil {
		if opts.Commitment != "" {
			obj["commitment"] = string(opts.Commitment)
		}
		if opts.MinContextSlot != nil {
			obj["minContextSlot"] = *opts.MinContextSlot
		}
	}

	params := []interface{}{blockHash}
	if len(obj) > 0 {
		params = append(params, obj)
	}

	err = cl.rpcClient.CallForInto(ctx, &out, "isBlockhashValid", params)
	return
}
