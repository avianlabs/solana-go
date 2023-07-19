package token

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeApprove             solana.TypedInstructionDecoder[*Approve]             = decode[*Approve]()
	DecodeApproveChecked      solana.TypedInstructionDecoder[*ApproveChecked]      = decode[*ApproveChecked]()
	DecodeBurn                solana.TypedInstructionDecoder[*Burn]                = decode[*Burn]()
	DecodeBurnChecked         solana.TypedInstructionDecoder[*BurnChecked]         = decode[*BurnChecked]()
	DecodeCloseAccount        solana.TypedInstructionDecoder[*CloseAccount]        = decode[*CloseAccount]()
	DecodeFreezeAccount       solana.TypedInstructionDecoder[*FreezeAccount]       = decode[*FreezeAccount]()
	DecodeInitializeAccount   solana.TypedInstructionDecoder[*InitializeAccount]   = decode[*InitializeAccount]()
	DecodeInitializeAccount2  solana.TypedInstructionDecoder[*InitializeAccount2]  = decode[*InitializeAccount2]()
	DecodeInitializeAccount3  solana.TypedInstructionDecoder[*InitializeAccount3]  = decode[*InitializeAccount3]()
	DecodeInitializeMint      solana.TypedInstructionDecoder[*InitializeMint]      = decode[*InitializeMint]()
	DecodeInitializeMint2     solana.TypedInstructionDecoder[*InitializeMint2]     = decode[*InitializeMint2]()
	DecodeInitializeMultisig  solana.TypedInstructionDecoder[*InitializeMultisig]  = decode[*InitializeMultisig]()
	DecodeInitializeMultisig2 solana.TypedInstructionDecoder[*InitializeMultisig2] = decode[*InitializeMultisig2]()
	DecodeMintTo              solana.TypedInstructionDecoder[*MintTo]              = decode[*MintTo]()
	DecodeMintToChecked       solana.TypedInstructionDecoder[*MintToChecked]       = decode[*MintToChecked]()
	DecodeRevoke              solana.TypedInstructionDecoder[*Revoke]              = decode[*Revoke]()
	DecodeSetAuthority        solana.TypedInstructionDecoder[*SetAuthority]        = decode[*SetAuthority]()
	DecodeSyncNative          solana.TypedInstructionDecoder[*SyncNative]          = decode[*SyncNative]()
	DecodeThawAccount         solana.TypedInstructionDecoder[*ThawAccount]         = decode[*ThawAccount]()
	DecodeTransfer            solana.TypedInstructionDecoder[*Transfer]            = decode[*Transfer]()
	DecodeTransferChecked     solana.TypedInstructionDecoder[*TransferChecked]     = decode[*TransferChecked]()
)

func decode[T any]() solana.TypedInstructionDecoder[T] {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
