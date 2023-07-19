package vote

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeAuthorize         = decode[*Authorize]()
	DecodeInitializeAccount = decode[*InitializeAccount]()
	DecodeWithdraw          = decode[*Withdraw]()
)

func decode[T any]() func(*solana.Message, int) (T, error) {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
