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

// buildCreateWithoutRent builds a Create instruction in the modern format
// which omits the rent sysvar account (as emitted by newer client libraries
// such as solana-kotlin).
func buildCreateWithoutRent(
	t *testing.T,
	payer, wallet, mint, tokenProgramID solana.PublicKey,
	idempotent bool,
) *Instruction {
	t.Helper()
	ata, _, err := solana.FindAssociatedTokenAddress(wallet, mint, tokenProgramID)
	require.NoError(t, err)
	create := Create{
		Payer:          payer,
		Wallet:         wallet,
		Mint:           mint,
		TokenProgramID: tokenProgramID,
		Idempotent:     idempotent,
	}
	create.AccountMetaSlice = solana.AccountMetaSlice{
		solana.NewAccountMeta(payer, true, true),
		solana.NewAccountMeta(ata, true, false),
		solana.NewAccountMeta(wallet, false, false),
		solana.NewAccountMeta(mint, false, false),
		solana.NewAccountMeta(solana.SystemProgramID, false, false),
		solana.NewAccountMeta(tokenProgramID, false, false),
	}
	return &Instruction{BaseVariant: bin.BaseVariant{
		Impl:   create,
		TypeID: bin.NoTypeIDDefaultID,
	}}
}

// testCreateKeys returns a deterministic set of keys for equivalence tests.
func testCreateKeys(t *testing.T) (payer, wallet, mint, tokenProgramID solana.PublicKey) {
	t.Helper()
	payerKey, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	payer = payerKey.PublicKey()
	wallet = payer
	mint = solana.MPK("G8iheDY9bGix5qCXEitCExLcgZzZrEemngk9cbTR3CQs")
	tokenProgramID = solana.TokenProgramID
	return payer, wallet, mint, tokenProgramID
}

func TestCreateAssertEquivalentOptionalRent(t *testing.T) {
	t.Parallel()
	payer, wallet, mint, tokenProgramID := testCreateKeys(t)

	withRent := NewCreateInstruction(payer, wallet, mint, tokenProgramID, false).Build()
	withoutRent := buildCreateWithoutRent(t, payer, wallet, mint, tokenProgramID, false)

	withRentImpl := withRent.Impl.(Create)
	withoutRentImpl := withoutRent.Impl.(Create)

	// Rent present on either side (or both, or neither) is equivalent.
	assert.NoError(t, withRentImpl.AssertEquivalent(&withoutRentImpl))
	assert.NoError(t, withoutRentImpl.AssertEquivalent(&withRentImpl))
	assert.NoError(t, withRentImpl.AssertEquivalent(&withRentImpl))
	assert.NoError(t, withoutRentImpl.AssertEquivalent(&withoutRentImpl))
}

func TestCreateAssertEquivalentRejectsBadTrailingAccount(t *testing.T) {
	t.Parallel()
	payer, wallet, mint, tokenProgramID := testCreateKeys(t)

	expected := NewCreateInstruction(payer, wallet, mint, tokenProgramID, false).Build().Impl.(Create)

	// A seventh account that is not the rent sysvar must be rejected.
	notRent := buildCreateWithoutRent(t, payer, wallet, mint, tokenProgramID, false).Impl.(Create)
	notRent.AccountMetaSlice = append(notRent.AccountMetaSlice,
		solana.NewAccountMeta(solana.SysVarClockPubkey, false, false))
	assert.Error(t, expected.AssertEquivalent(&notRent))

	// A writable rent sysvar must be rejected.
	writableRent := buildCreateWithoutRent(t, payer, wallet, mint, tokenProgramID, false).Impl.(Create)
	writableRent.AccountMetaSlice = append(writableRent.AccountMetaSlice,
		solana.NewAccountMeta(solana.SysVarRentPubkey, true, false))
	assert.Error(t, expected.AssertEquivalent(&writableRent))

	// More than seven accounts must be rejected.
	extra := NewCreateInstruction(payer, wallet, mint, tokenProgramID, false).Build().Impl.(Create)
	extra.AccountMetaSlice = append(extra.AccountMetaSlice,
		solana.NewAccountMeta(solana.SysVarClockPubkey, false, false))
	assert.Error(t, expected.AssertEquivalent(&extra))

	// A missing required account must still be rejected.
	truncated := buildCreateWithoutRent(t, payer, wallet, mint, tokenProgramID, false).Impl.(Create)
	truncated.AccountMetaSlice = truncated.AccountMetaSlice[:5]
	assert.Error(t, expected.AssertEquivalent(&truncated))

	// A differing required account must still be rejected, with or without rent.
	otherWallet, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	differing := buildCreateWithoutRent(t, payer, otherWallet.PublicKey(), mint, tokenProgramID, false).Impl.(Create)
	assert.Error(t, expected.AssertEquivalent(&differing))
}

func TestMessageAssertEquivalentToleratesMissingRent(t *testing.T) {
	t.Parallel()
	payer, wallet, mint, tokenProgramID := testCreateKeys(t)
	blockhash := solana.MustHashFromBase58("AnL7vVGidfdxZBFqkxLvVwXgW9pDSZC5kK6d5kyRaXEa")

	expected, err := solana.NewTransaction(
		[]solana.Instruction{
			NewCreateInstruction(payer, wallet, mint, tokenProgramID, false).Build(),
		},
		blockhash,
		solana.TransactionPayer(payer),
	)
	require.NoError(t, err)

	// Simulate a client-built transaction in the modern rent-less format
	// arriving over the wire.
	clientTx, err := solana.NewTransaction(
		[]solana.Instruction{
			buildCreateWithoutRent(t, payer, wallet, mint, tokenProgramID, false),
		},
		blockhash,
		solana.TransactionPayer(payer),
	)
	require.NoError(t, err)
	wire, err := clientTx.MarshalBinary()
	require.NoError(t, err)
	received := &solana.Transaction{}
	require.NoError(t, received.UnmarshalWithDecoder(bin.NewCompactU16Decoder(wire)))

	assert.NoError(t, expected.Message.AssertEquivalent(received.Message))
	assert.NoError(t, received.Message.AssertEquivalent(expected.Message))
}

func TestMessageAssertEquivalentStillRejectsOtherDifferences(t *testing.T) {
	t.Parallel()
	payer, wallet, mint, tokenProgramID := testCreateKeys(t)
	blockhash := solana.MustHashFromBase58("AnL7vVGidfdxZBFqkxLvVwXgW9pDSZC5kK6d5kyRaXEa")

	expected, err := solana.NewTransaction(
		[]solana.Instruction{
			NewCreateInstruction(payer, wallet, mint, tokenProgramID, false).Build(),
		},
		blockhash,
		solana.TransactionPayer(payer),
	)
	require.NoError(t, err)

	// A rent-less transaction for a different wallet must still fail even
	// though the rent difference itself is tolerated.
	otherWallet, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	differing, err := solana.NewTransaction(
		[]solana.Instruction{
			buildCreateWithoutRent(t, payer, otherWallet.PublicKey(), mint, tokenProgramID, false),
		},
		blockhash,
		solana.TransactionPayer(payer),
	)
	require.NoError(t, err)
	assert.Error(t, expected.Message.AssertEquivalent(differing.Message))

	// A different blockhash must still fail.
	otherBlockhash, err := solana.NewTransaction(
		[]solana.Instruction{
			buildCreateWithoutRent(t, payer, wallet, mint, tokenProgramID, false),
		},
		solana.MustHashFromBase58("EkSnNWid2cvwEVnVx9aBqawnmiCNiDgp3gUdkDPTKN1N"),
		solana.TransactionPayer(payer),
	)
	require.NoError(t, err)
	assert.Error(t, expected.Message.AssertEquivalent(otherBlockhash.Message))
}
