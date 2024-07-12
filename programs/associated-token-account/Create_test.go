package associatedtokenaccount

import (
	"testing"

	bin "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/require"
)

func TestCreateNonIdempotentData(t *testing.T) {
	t.Parallel()
	wallet, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	c := NewCreateInstruction(
		wallet.PublicKey(),
		wallet.PublicKey(),
		solana.MPK("G8iheDY9bGix5qCXEitCExLcgZzZrEemngk9cbTR3CQs"),
		solana.MustPublicKeyFromBase58("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"),
		false,
	)

	data, err := c.Build().Data()
	require.NoError(t, err)

	assert.Len(t, data, 0)
}

func TestCreateIdempotentData(t *testing.T) {
	t.Parallel()
	wallet, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	c := NewCreateInstruction(
		wallet.PublicKey(),
		wallet.PublicKey(),
		solana.MPK("G8iheDY9bGix5qCXEitCExLcgZzZrEemngk9cbTR3CQs"),
		solana.MustPublicKeyFromBase58("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"),
		true,
	)

	data, err := c.Build().Data()
	require.NoError(t, err)

	assert.Equal(t, []byte{1}, data)
}

func TestEncodeRoundtrip(t *testing.T) {
	t.Parallel()
	wallet, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	c := NewCreateInstruction(
		wallet.PublicKey(),
		wallet.PublicKey(),
		solana.MPK("G8iheDY9bGix5qCXEitCExLcgZzZrEemngk9cbTR3CQs"),
		solana.MustPublicKeyFromBase58("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"),
		true,
	).Build()

	blockhash := solana.MustHashFromBase58("AnL7vVGidfdxZBFqkxLvVwXgW9pDSZC5kK6d5kyRaXEa")
	tx, err := solana.NewTransaction([]solana.Instruction{c}, blockhash)
	require.NoError(t, err)

	data, err := tx.MarshalBinary()
	require.NoError(t, err)

	decoded := &solana.Transaction{}
	err = decoded.UnmarshalWithDecoder(bin.NewCompactU16Decoder(data))
	require.NoError(t, err)

	err = tx.Message.AssertEquivalent(decoded.Message)
	require.NoError(t, err)
}
