package instruction

import (
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
)

type Instruction struct {
	solana.Instruction
}

func (v Instruction) AsTokenTransfer() (*token.Transfer, bool) {
	ti, ok := v.Instruction.(*token.Instruction)
	if !ok {
		return nil, false
	}
	tt, ok := ti.BaseVariant.Impl.(*token.Transfer)
	if !ok {
		return nil, false
	}
	return tt, true
}

func (v Instruction) AsTokenTransferChecked() (*token.TransferChecked, bool) {
	ti, ok := v.Instruction.(*token.Instruction)
	if !ok {
		return nil, false
	}
	ttc, ok := ti.BaseVariant.Impl.(*token.TransferChecked)
	if !ok {
		return nil, false
	}
	return ttc, true
}

func (v Instruction) AsSystemTransfer() (*system.Transfer, bool) {
	si, ok := v.Instruction.(*system.Instruction)
	if !ok {
		return nil, false
	}
	st, ok := si.BaseVariant.Impl.(*system.Transfer)
	if !ok {
		return nil, false
	}
	return st, true
}

func (v Instruction) AsCreateAssociatedTokenAccount() (*associatedtokenaccount.Create, bool) {
	ai, ok := v.Instruction.(*associatedtokenaccount.Instruction)
	if !ok {
		return nil, false
	}
	ci, ok := ai.BaseVariant.Impl.(*associatedtokenaccount.Create)
	if !ok {
		return nil, false
	}
	return ci, true
}

func (v Instruction) AsCloseTokenAccount() (*token.CloseAccount, bool) {
	ti, ok := v.Instruction.(*token.Instruction)
	if !ok {
		return nil, false
	}
	ta, ok := ti.BaseVariant.Impl.(*token.CloseAccount)
	if !ok {
		return nil, false
	}
	return ta, true
}
