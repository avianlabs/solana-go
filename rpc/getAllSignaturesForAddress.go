package rpc

import (
	"context"

	"github.com/gagliardetto/solana-go"
)

type GetAllSignaturesForAddressOpts struct {
	// (optional) PerRequestLimit is the limit on number of signatures to fetch
	// 'per-request', the default (and max) if none provided is 1000.
	PerRequestLimit *int

	// (optional) Start searching backwards from this transaction signature. If
	// not provided the search starts from the top of the highest max confirmed
	// block.
	Before solana.Signature

	// (optional) Search until this transaction signature, if found before
	// limit reached.
	Until solana.Signature

	// (optional) Commitment; "processed" is not supported.
	// If parameter not provided, the default is "finalized".
	Commitment CommitmentType
}

// GetAllSignaturesForAddressWithOpts is the same as
// 'GetSignaturesForAddressWithOpts' except it will continue requesting
// until _all_ signatures have been fetched (within the bounds) - rather than
// stopping after 'limit'.
func (v *Client) GetAllSignaturesForAddressWithOpts(ctx context.Context, addr solana.PublicKey, opts *GetAllSignaturesForAddressOpts) ([]*TransactionSignature, error) {
	var (
		limit      = 1000
		before     solana.Signature
		until      solana.Signature
		commitment = CommitmentFinalized
		all        []*TransactionSignature
	)
	if opts != nil {
		before = opts.Before
		until = opts.Until
		if opts.PerRequestLimit != nil {
			limit = *opts.PerRequestLimit
		}
		if opts.Commitment != "" {
			commitment = opts.Commitment
		}
	}
	for {
		res, err := v.GetSignaturesForAddressWithOpts(ctx, addr, &GetSignaturesForAddressOpts{
			Limit:      &limit,
			Until:      until,
			Before:     before,
			Commitment: commitment,
		})
		if err != nil {
			return nil, err
		}
		all = append(all, res...)
		if len(res) < limit {
			// There's no more, either we've reached 'until' or the beginning
			// of Txs for this addr.
			break
		}
		// Go again but move the search window.
		before = res[len(res)-1].Signature
	}
	return all, nil
}
