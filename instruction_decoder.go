package solana

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
)

// TypedInstructionDecoder implementations decode the instruction at the
// provided index in the message to the provided type.
type TypedInstructionDecoder[A any] func(*Message, int) (A, error)

// DecodeInstructionType decodes instruction at index in the message using the
// provided decoders. Returning the requested type. It's expected that this is
// used in conjunction with typed instructions in the program packages.
func DecodeInstructionType[A interface {
	Instruction
	Obtain(*bin.VariantDefinition) (bin.TypeID, string, interface{})
}, B any](
	expectedProgramID PublicKey,
	def *bin.VariantDefinition,
	decodeFn func([]*AccountMeta, []byte) (A, error),
) TypedInstructionDecoder[B] {
	var (
		a A
		b B
	)
	return func(msg *Message, index int) (B, error) {
		if len(msg.Instructions) <= index {
			return b, fmt.Errorf("transaction doesn't have an instruction at index '%d'", index)
		}
		instruction := msg.Instructions[index]
		accs, err := instruction.ResolveInstructionAccounts(msg)
		if err != nil {
			return b, fmt.Errorf("instruction '%d': failed to resolve accounts: %w", index, err)
		}
		programID, err := msg.ResolveProgramIDIndex(instruction.ProgramIDIndex)
		if err != nil {
			return b, fmt.Errorf("instruction '%d': failed to resolve program ID: %w", index, err)
		}
		if !programID.Equals(expectedProgramID) {
			return b, fmt.Errorf("instruction '%d': programID (%s) doesn't match expected value '%s'", index, programID, expectedProgramID)
		}
		decoded, err := decodeFn(accs, instruction.Data)
		if err != nil {
			return b, fmt.Errorf("instruction '%d': failed to decode as '%T': %w", index, a, err)
		}
		_, _, obtained := decoded.Obtain(def)
		v, ok := obtained.(B)
		if !ok {
			return b, fmt.Errorf("instruction '%d': obtained type '%T' doesn't match expected type '%T'", index, obtained, b)
		}
		return v, nil
	}
}
