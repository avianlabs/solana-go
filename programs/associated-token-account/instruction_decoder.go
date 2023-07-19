package associatedtokenaccount

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeCreate solana.TypedInstructionDecoder[*Create] = decode[*Create]()
)

func decode[T any]() solana.TypedInstructionDecoder[T] {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
