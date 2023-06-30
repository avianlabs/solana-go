package instruction

import (
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/programs/vote"
)

// Decoder makes a best effort to decode the transaction instruction based on
// the system program. If unable to decode an instruction, it's left as a
// CompiledInstruction.
type Decoder struct {
	tx *solana.Transaction
}

func NewDecoder(tx *solana.Transaction) Decoder {
	return Decoder{tx}
}

// DecodeAll instructions in the transaction.
func (v Decoder) DecodeAll() ([]*Instruction, error) {
	ins := make([]*Instruction, 0, len(v.tx.Message.Instructions))
	for i := range v.tx.Message.Instructions {
		in, err := v.Decode(i)
		if err != nil {
			return nil, err
		}
		ins = append(ins, in)
	}
	return ins, nil
}

// Decode instruction at index i.
func (v Decoder) Decode(i int) (*Instruction, error) {
	in := v.tx.Message.Instructions[i]
	accs, err := in.ResolveInstructionAccounts(&v.tx.Message)
	if err != nil {
		return nil, err
	}
	progID, err := v.tx.ResolveProgramIDIndex(in.ProgramIDIndex)
	if err != nil {
		return nil, err
	}
	var res solana.Instruction
	switch progID {
	case token.ProgramID:
		res, err = token.DecodeInstruction(accs, in.Data)
	case system.ProgramID:
		res, err = system.DecodeInstruction(accs, in.Data)
	case vote.ProgramID:
		res, err = vote.DecodeInstruction(accs, in.Data)
	case associatedtokenaccount.ProgramID:
		res, err = associatedtokenaccount.DecodeInstruction(accs, in.Data)
	default:
		res = NewCompiled(in, accs, progID)
	}
	if err != nil {
		return nil, err
	}
	return &Instruction{res}, nil
}
