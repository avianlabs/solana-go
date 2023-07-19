package vote

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeAuthorize         solana.TypedInstructionDecoder[*Authorize]         = decode[*Authorize]()
	DecodeInitializeAccount solana.TypedInstructionDecoder[*InitializeAccount] = decode[*InitializeAccount]()
	DecodeWithdraw          solana.TypedInstructionDecoder[*Withdraw]          = decode[*Withdraw]()
)

func decode[T any]() solana.TypedInstructionDecoder[T] {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
