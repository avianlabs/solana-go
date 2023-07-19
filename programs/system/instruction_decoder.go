package system

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeAdvanceNonceAccount    = decode[*AdvanceNonceAccount]()
	DecodeAllocate               = decode[*Allocate]()
	DecodeAllocateWithSeed       = decode[*AllocateWithSeed]()
	DecodeAssign                 = decode[*Assign]()
	DecodeAssignWithSeed         = decode[*AssignWithSeed]()
	DecodeAuthorizeNonceAccount  = decode[*AuthorizeNonceAccount]()
	DecodeCreateAccount          = decode[*CreateAccount]()
	DecodeCreateAccountWithSeed  = decode[*CreateAccountWithSeed]()
	DecodeInitializeNonceAccount = decode[*InitializeNonceAccount]()
	DecodeTransfer               = decode[*Transfer]()
	DecodeTransferWithSeed       = decode[*TransferWithSeed]()
	DecodeWithdrawNonceAccount   = decode[*WithdrawNonceAccount]()
)

func decode[T any]() func(*solana.Message, int) (T, error) {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
