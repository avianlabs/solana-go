package system

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeAdvanceNonceAccount    solana.TypedInstructionDecoder[*AdvanceNonceAccount]    = decode[*AdvanceNonceAccount]()
	DecodeAllocate               solana.TypedInstructionDecoder[*Allocate]               = decode[*Allocate]()
	DecodeAllocateWithSeed       solana.TypedInstructionDecoder[*AllocateWithSeed]       = decode[*AllocateWithSeed]()
	DecodeAssign                 solana.TypedInstructionDecoder[*Assign]                 = decode[*Assign]()
	DecodeAssignWithSeed         solana.TypedInstructionDecoder[*AssignWithSeed]         = decode[*AssignWithSeed]()
	DecodeAuthorizeNonceAccount  solana.TypedInstructionDecoder[*AuthorizeNonceAccount]  = decode[*AuthorizeNonceAccount]()
	DecodeCreateAccount          solana.TypedInstructionDecoder[*CreateAccount]          = decode[*CreateAccount]()
	DecodeCreateAccountWithSeed  solana.TypedInstructionDecoder[*CreateAccountWithSeed]  = decode[*CreateAccountWithSeed]()
	DecodeInitializeNonceAccount solana.TypedInstructionDecoder[*InitializeNonceAccount] = decode[*InitializeNonceAccount]()
	DecodeTransfer               solana.TypedInstructionDecoder[*Transfer]               = decode[*Transfer]()
	DecodeTransferWithSeed       solana.TypedInstructionDecoder[*TransferWithSeed]       = decode[*TransferWithSeed]()
	DecodeWithdrawNonceAccount   solana.TypedInstructionDecoder[*WithdrawNonceAccount]   = decode[*WithdrawNonceAccount]()
)

func decode[T any]() solana.TypedInstructionDecoder[T] {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
