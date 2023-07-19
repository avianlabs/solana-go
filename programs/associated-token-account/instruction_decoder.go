package associatedtokenaccount

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeCreate = decode[*Create]()
)

func decode[T any]() func(*solana.Message, int) (T, error) {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
