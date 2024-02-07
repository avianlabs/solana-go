package token2022

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeApprove             = decode[*Approve]()
	DecodeApproveChecked      = decode[*ApproveChecked]()
	DecodeBurn                = decode[*Burn]()
	DecodeBurnChecked         = decode[*BurnChecked]()
	DecodeCloseAccount        = decode[*CloseAccount]()
	DecodeFreezeAccount       = decode[*FreezeAccount]()
	DecodeInitializeAccount   = decode[*InitializeAccount]()
	DecodeInitializeAccount2  = decode[*InitializeAccount2]()
	DecodeInitializeAccount3  = decode[*InitializeAccount3]()
	DecodeInitializeMint      = decode[*InitializeMint]()
	DecodeInitializeMint2     = decode[*InitializeMint2]()
	DecodeInitializeMultisig  = decode[*InitializeMultisig]()
	DecodeInitializeMultisig2 = decode[*InitializeMultisig2]()
	DecodeMintTo              = decode[*MintTo]()
	DecodeMintToChecked       = decode[*MintToChecked]()
	DecodeRevoke              = decode[*Revoke]()
	DecodeSetAuthority        = decode[*SetAuthority]()
	DecodeSyncNative          = decode[*SyncNative]()
	DecodeThawAccount         = decode[*ThawAccount]()
	DecodeTransfer            = decode[*Transfer]()
	DecodeTransferChecked     = decode[*TransferChecked]()
)

func decode[T any]() func(*solana.Message, int) (T, error) {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
