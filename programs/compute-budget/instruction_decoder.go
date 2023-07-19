package computebudget

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeRequestHeapFrame       solana.TypedInstructionDecoder[*RequestHeapFrame]       = decode[*RequestHeapFrame]()
	DecodeRequestUnitsDeprecated solana.TypedInstructionDecoder[*RequestUnitsDeprecated] = decode[*RequestUnitsDeprecated]()
	DecodeSetComputeUnitLimit    solana.TypedInstructionDecoder[*SetComputeUnitLimit]    = decode[*SetComputeUnitLimit]()
	DecodeSetComputeUnitPrice    solana.TypedInstructionDecoder[*SetComputeUnitPrice]    = decode[*SetComputeUnitPrice]()
)

func decode[T any]() solana.TypedInstructionDecoder[T] {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
