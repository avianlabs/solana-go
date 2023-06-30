package instruction

import "github.com/gagliardetto/solana-go"

type Compiled struct {
	accs   []*solana.AccountMeta
	progID solana.PublicKey
	solana.CompiledInstruction
}

var _ solana.Instruction = &Compiled{}

func NewCompiled(ci solana.CompiledInstruction, accs []*solana.AccountMeta, progID solana.PublicKey) *Compiled {
	return &Compiled{
		CompiledInstruction: ci,
		progID:              progID,
		accs:                accs,
	}
}

func (v *Compiled) ProgramID() solana.PublicKey {
	return v.progID
}

func (v *Compiled) Accounts() []*solana.AccountMeta {
	return v.accs
}

func (v *Compiled) Data() ([]byte, error) {
	bytes, err := v.CompiledInstruction.Data.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
