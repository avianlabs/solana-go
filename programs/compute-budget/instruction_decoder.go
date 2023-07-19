package computebudget

import (
	"github.com/gagliardetto/solana-go"
)

var (
	DecodeRequestHeapFrame       = decode[*RequestHeapFrame]()
	DecodeRequestUnitsDeprecated = decode[*RequestUnitsDeprecated]()
	DecodeSetComputeUnitLimit    = decode[*SetComputeUnitLimit]()
	DecodeSetComputeUnitPrice    = decode[*SetComputeUnitPrice]()
)

func decode[T any]() func(*solana.Message, int) (T, error) {
	return solana.DecodeInstructionType[*Instruction, T](
		ProgramID,
		InstructionImplDef,
		DecodeInstruction,
	)
}
