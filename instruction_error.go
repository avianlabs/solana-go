package solana

import (
	"fmt"
	"sync"
)

var (
	customErrorResolvers  = map[PublicKey]func(int) (error, bool){}
	customErrorResolverMU sync.RWMutex
)

func RegisterCustomInstructionErrorResolver(progID PublicKey, f func(int) (error, bool)) {
	customErrorResolverMU.Lock()
	defer customErrorResolverMU.Unlock()
	customErrorResolvers[progID] = f
}

func ResolveCustomInstructionError(progID PublicKey, code int) (error, bool) {
	customErrorResolverMU.RLock()
	defer customErrorResolverMU.RUnlock()
	resolver, ok := customErrorResolvers[progID]
	if !ok {
		return nil, false
	}
	return resolver(code)
}

type InstructionError interface {
	error
	xx_isInstructionError()
}

func CustomInstructionErrorResolver(tx *Transaction, index uint16) func(code int) (error, bool) {
	return func(code int) (error, bool) {
		in := tx.Message.Instructions[index]
		prog, err := tx.ResolveProgramIDIndex(in.ProgramIDIndex)
		if err != nil {
			return nil, false //nolint: nilerr
		}
		return ResolveCustomInstructionError(prog, code)
	}
}

func ParseInstructionError(err interface{}, resolve func(int) (error, bool)) (InstructionError, bool) {
	switch t := err.(type) {
	case string:
		return parseInstructionErrorString(t)
	case map[string]interface{}:
		return parseInstructionErrorObject(t, resolve)
	default:
		return nil, false
	}
}

func parseInstructionErrorObject(err map[string]interface{}, resolve func(code int) (error, bool)) (InstructionError, bool) {
	code, ok := asFloat64(err["Custom"])
	if !ok {
		return nil, false
	}
	cause, _ := resolve(int(code))
	return &InstructionError_Custom{
		Code:  uint32(code),
		Cause: cause,
	}, true
}

func parseInstructionErrorString(err string) (InstructionError, bool) {
	switch err {
	case "GenericError":
		return InstructionError_GenericError{}, true
	case "InvalidArgument":
		return InstructionError_InvalidArgument{}, true
	case "InvalidInstructionData":
		return InstructionError_InvalidInstructionData{}, true
	case "InvalidAccountData":
		return InstructionError_InvalidAccountData{}, true
	case "AccountDataTooSmall":
		return InstructionError_AccountDataTooSmall{}, true
	case "InsufficientFunds":
		return InstructionError_InsufficientFunds{}, true
	case "IncorrectProgramId":
		return InstructionError_IncorrectProgramId{}, true
	case "MissingRequiredSignature":
		return InstructionError_MissingRequiredSignature{}, true
	case "AccountAlreadyInitialized":
		return InstructionError_AccountAlreadyInitialized{}, true
	case "UninitializedAccount":
		return InstructionError_UninitializedAccount{}, true
	case "UnbalancedInstruction":
		return InstructionError_UnbalancedInstruction{}, true
	case "ModifiedProgramId":
		return InstructionError_ModifiedProgramId{}, true
	case "ExternalAccountLamportSpend":
		return InstructionError_ExternalAccountLamportSpend{}, true
	case "ExternalAccountDataModified":
		return InstructionError_ExternalAccountDataModified{}, true
	case "ReadonlyLamportChange":
		return InstructionError_ReadonlyLamportChange{}, true
	case "ReadonlyDataModified":
		return InstructionError_ReadonlyDataModified{}, true
	case "DuplicateAccountIndex":
		return InstructionError_DuplicateAccountIndex{}, true
	case "ExecutableModified":
		return InstructionError_ExecutableModified{}, true
	case "RentEpochModified":
		return InstructionError_RentEpochModified{}, true
	case "NotEnoughAccountKeys":
		return InstructionError_NotEnoughAccountKeys{}, true
	case "AccountDataSizeChanged":
		return InstructionError_AccountDataSizeChanged{}, true
	case "AccountNotExecutable":
		return InstructionError_AccountNotExecutable{}, true
	case "AccountBorrowFailed":
		return InstructionError_AccountBorrowFailed{}, true
	case "AccountBorrowOutstanding":
		return InstructionError_AccountBorrowOutstanding{}, true
	case "DuplicateAccountOutOfSync":
		return InstructionError_DuplicateAccountOutOfSync{}, true
	case "InvalidError":
		return InstructionError_InvalidError{}, true
	case "ExecutableDataModified":
		return InstructionError_ExecutableDataModified{}, true
	case "ExecutableLamportChange":
		return InstructionError_ExecutableLamportChange{}, true
	case "ExecutableAccountNotRentExempt":
		return InstructionError_ExecutableAccountNotRentExempt{}, true
	case "UnsupportedProgramId":
		return InstructionError_UnsupportedProgramId{}, true
	case "CallDepth":
		return InstructionError_CallDepth{}, true
	case "MissingAccount":
		return InstructionError_MissingAccount{}, true
	case "ReentrancyNotAllowed":
		return InstructionError_ReentrancyNotAllowed{}, true
	case "MaxSeedLengthExceeded":
		return InstructionError_MaxSeedLengthExceeded{}, true
	case "InvalidSeeds":
		return InstructionError_InvalidSeeds{}, true
	case "InvalidRealloc":
		return InstructionError_InvalidRealloc{}, true
	case "ComputationalBudgetExceeded":
		return InstructionError_ComputationalBudgetExceeded{}, true
	case "PrivilegeEscalation":
		return InstructionError_PrivilegeEscalation{}, true
	case "ProgramEnvironmentSetupFailure":
		return InstructionError_ProgramEnvironmentSetupFailure{}, true
	case "ProgramFailedToComplete":
		return InstructionError_ProgramFailedToComplete{}, true
	case "ProgramFailedToCompile":
		return InstructionError_ProgramFailedToCompile{}, true
	case "Immutable":
		return InstructionError_Immutable{}, true
	case "IncorrectAuthority":
		return InstructionError_IncorrectAuthority{}, true
	case "AccountNotRentExempt":
		return InstructionError_AccountNotRentExempt{}, true
	case "InvalidAccountOwner":
		return InstructionError_InvalidAccountOwner{}, true
	case "ArithmeticOverflow":
		return InstructionError_ArithmeticOverflow{}, true
	case "UnsupportedSysvar":
		return InstructionError_UnsupportedSysvar{}, true
	case "IllegalOwner":
		return InstructionError_IllegalOwner{}, true
	default:
		return InstructionError_Undefined(err), true
	}
}

type InstructionError_Undefined string

func (v InstructionError_Undefined) Error() string {
	return string(v)
}
func (InstructionError_Undefined) xx_isInstructionError() {}

// Defined [here](https://github.com/solana-labs/solana/blob/f6371cce176d481b4132e5061262ca015db0f8b1/sdk/program/src/instruction.rs#L23).

// The program instruction returned an error
type InstructionError_GenericError struct{}

func (InstructionError_GenericError) Error() string          { return "generic instruction error" }
func (InstructionError_GenericError) xx_isInstructionError() {}

// The arguments provided to a program were invalid
type InstructionError_InvalidArgument struct{}

func (InstructionError_InvalidArgument) Error() string          { return "invalid program argument" }
func (InstructionError_InvalidArgument) xx_isInstructionError() {}

// An instruction's data contents were invalid
type InstructionError_InvalidInstructionData struct{}

func (InstructionError_InvalidInstructionData) Error() string          { return "invalid instruction data" }
func (InstructionError_InvalidInstructionData) xx_isInstructionError() {}

// An account's data contents was invalid
type InstructionError_InvalidAccountData struct{}

func (InstructionError_InvalidAccountData) Error() string {
	return "invalid account data for instruction"
}
func (InstructionError_InvalidAccountData) xx_isInstructionError() {}

// An account's data was too small
type InstructionError_AccountDataTooSmall struct{}

func (InstructionError_AccountDataTooSmall) Error() string {
	return "account data too small for instruction"
}
func (InstructionError_AccountDataTooSmall) xx_isInstructionError() {}

// An account's balance was too small to complete the instruction
type InstructionError_InsufficientFunds struct{}

func (InstructionError_InsufficientFunds) Error() string          { return "insufficient funds for instruction" }
func (InstructionError_InsufficientFunds) xx_isInstructionError() {}

// The account did not have the expected program id
type InstructionError_IncorrectProgramId struct{}

func (InstructionError_IncorrectProgramId) Error() string {
	return "incorrect program id for instruction"
}
func (InstructionError_IncorrectProgramId) xx_isInstructionError() {}

// A signature was required but not found
type InstructionError_MissingRequiredSignature struct{}

func (InstructionError_MissingRequiredSignature) Error() string {
	return "missing required signature for instruction"
}
func (InstructionError_MissingRequiredSignature) xx_isInstructionError() {}

// An initialize instruction was sent to an account that has already been initialized.
type InstructionError_AccountAlreadyInitialized struct{}

func (InstructionError_AccountAlreadyInitialized) Error() string {
	return "instruction requires an uninitialized account"
}
func (InstructionError_AccountAlreadyInitialized) xx_isInstructionError() {}

// An attempt to operate on an account that hasn't been initialized.
type InstructionError_UninitializedAccount struct{}

func (InstructionError_UninitializedAccount) Error() string {
	return "instruction requires an initialized account"
}
func (InstructionError_UninitializedAccount) xx_isInstructionError() {}

// UnbalancedInstruction is an error that occurs when a program's instruction
// lamport balance does not equal the balance after the instruction.
type InstructionError_UnbalancedInstruction struct{}

func (InstructionError_UnbalancedInstruction) Error() string {
	return "sum of account balances before and after instruction do not match"
}
func (InstructionError_UnbalancedInstruction) xx_isInstructionError() {}

// ModifiedProgramId is an error that occurs when a program modifies an
// account's program id.
type InstructionError_ModifiedProgramId struct{}

func (InstructionError_ModifiedProgramId) Error() string {
	return "instruction modified the program id of an account"
}
func (InstructionError_ModifiedProgramId) xx_isInstructionError() {}

// ExternalAccountLamportSpend is an error that occurs when a program spends
// the lamports of an account that doesn't belong to it.
type InstructionError_ExternalAccountLamportSpend struct{}

func (InstructionError_ExternalAccountLamportSpend) Error() string {
	return "instruction spent from the balance of an account it does not own"
}
func (InstructionError_ExternalAccountLamportSpend) xx_isInstructionError() {}

// ExternalAccountDataModified is an error that occurs when a program modifies
// the data of an account that doesn't belong to it.
type InstructionError_ExternalAccountDataModified struct{}

func (InstructionError_ExternalAccountDataModified) Error() string {
	return "instruction modified data of an account it does not own"
}
func (InstructionError_ExternalAccountDataModified) xx_isInstructionError() {}

// ReadonlyLamportChange is an error that occurs when a read-only account's
// lamports are modified.
type InstructionError_ReadonlyLamportChange struct{}

func (InstructionError_ReadonlyLamportChange) Error() string {
	return "instruction changed the balance of a read-only account"
}
func (InstructionError_ReadonlyLamportChange) xx_isInstructionError() {}

// ReadonlyDataModified is an error that occurs when a read-only account's data
// is modified.
type InstructionError_ReadonlyDataModified struct{}

func (InstructionError_ReadonlyDataModified) Error() string {
	return "instruction modified data of a read-only account"
}
func (InstructionError_ReadonlyDataModified) xx_isInstructionError() {}

// DuplicateAccountIndex is an error that occurs when an account is referenced
// more than once in a single instruction.
type InstructionError_DuplicateAccountIndex struct{}

func (InstructionError_DuplicateAccountIndex) Error() string {
	return "instruction contains duplicate accounts"
}
func (InstructionError_DuplicateAccountIndex) xx_isInstructionError() {}

// ExecutableModified is an error that occurs when an instruction changes the
// executable bit of an account.
type InstructionError_ExecutableModified struct{}

func (InstructionError_ExecutableModified) Error() string {
	return "instruction changed executable bit of an account"
}

func (InstructionError_ExecutableModified) xx_isInstructionError() {}

// RentEpochModified is an error that occurs when an instruction modifies the
// rent epoch of an account.
type InstructionError_RentEpochModified struct{}

func (InstructionError_RentEpochModified) Error() string {
	return "instruction modified rent epoch of an account"
}

func (InstructionError_RentEpochModified) xx_isInstructionError() {}

// NotEnoughAccountKeys is an error that occurs when an instruction does not
// have enough account keys.
type InstructionError_NotEnoughAccountKeys struct{}

func (InstructionError_NotEnoughAccountKeys) Error() string {
	return "insufficient account keys for instruction"
}

func (InstructionError_NotEnoughAccountKeys) xx_isInstructionError() {}

// AccountDataSizeChanged is an error that occurs when a non-system program
// changes the size of an account data.
type InstructionError_AccountDataSizeChanged struct{}

func (InstructionError_AccountDataSizeChanged) Error() string {
	return "non-system instruction changed account size"
}

func (InstructionError_AccountDataSizeChanged) xx_isInstructionError() {}

// AccountNotExecutable is an error that occurs when an instruction expects an
// executable account.
type InstructionError_AccountNotExecutable struct{}

func (InstructionError_AccountNotExecutable) Error() string {
	return "instruction expected an executable account"
}

func (InstructionError_AccountNotExecutable) xx_isInstructionError() {}

// AccountBorrowFailed is an error that occurs when an instruction fails to
// borrow a reference for an account.
type InstructionError_AccountBorrowFailed struct{}

func (InstructionError_AccountBorrowFailed) Error() string {
	return "instruction tries to borrow reference for an account which is already borrowed"
}

func (InstructionError_AccountBorrowFailed) xx_isInstructionError() {}

// AccountBorrowOutstanding is an error that occurs when an account data has an
// outstanding reference after a program's execution.
type InstructionError_AccountBorrowOutstanding struct{}

func (InstructionError_AccountBorrowOutstanding) Error() string {
	return "instruction left account with an outstanding borrowed reference"
}
func (InstructionError_AccountBorrowOutstanding) xx_isInstructionError() {}

// The same account was multiply passed to an on-chain program's entrypoint,
// but the program modified them differently.  A program can only modify one
// instance of the account because the runtime cannot determine which changes
// to pick or how to merge them if both are modified
type InstructionError_DuplicateAccountOutOfSync struct{}

func (InstructionError_DuplicateAccountOutOfSync) Error() string {
	return "instruction modifications of multiply-passed account differ"
}
func (InstructionError_DuplicateAccountOutOfSync) xx_isInstructionError() {}

// Allows on-chain programs to implement program-specific error types and see
// them returned by the Solana runtime. A program-specific error may be any
// type that is represented as or serialized to a u32 integer.
type InstructionError_Custom struct {
	Code  uint32
	Cause error
}

func (v *InstructionError_Custom) Error() string {
	if v.Cause != nil {
		return v.Cause.Error()
	}
	return fmt.Sprintf("custom program error: %#x", v.Code)
}
func (v *InstructionError_Custom) Unwrap() error        { return v.Cause }
func (*InstructionError_Custom) xx_isInstructionError() {}

// The return value from the program was invalid.  Valid errors are either a
// defined builtin error value or a user-defined error in the lower 32 bits.
type InstructionError_InvalidError struct{}

func (InstructionError_InvalidError) Error() string          { return "program returned invalid error code" }
func (InstructionError_InvalidError) xx_isInstructionError() {}

// Executable account's data was modified
type InstructionError_ExecutableDataModified struct{}

func (InstructionError_ExecutableDataModified) Error() string {
	return "instruction changed executable accounts data"
}
func (InstructionError_ExecutableDataModified) xx_isInstructionError() {}

// Executable account's lamports modified
type InstructionError_ExecutableLamportChange struct{}

func (InstructionError_ExecutableLamportChange) Error() string {
	return "instruction changed the balance of a executable account"
}
func (InstructionError_ExecutableLamportChange) xx_isInstructionError() {}

// Executable accounts must be rent exempt
type InstructionError_ExecutableAccountNotRentExempt struct{}

func (InstructionError_ExecutableAccountNotRentExempt) Error() string {
	return "executable accounts must be rent exempt"
}
func (InstructionError_ExecutableAccountNotRentExempt) xx_isInstructionError() {}

// Unsupported program id
type InstructionError_UnsupportedProgramId struct{}

func (InstructionError_UnsupportedProgramId) Error() string          { return "Unsupported program id" }
func (InstructionError_UnsupportedProgramId) xx_isInstructionError() {}

// Cross-program invocation call depth too deep
type InstructionError_CallDepth struct{}

func (InstructionError_CallDepth) Error() string {
	return "Cross-program invocation call depth too deep"
}
func (InstructionError_CallDepth) xx_isInstructionError() {}

// An account required by the instruction is missing
type InstructionError_MissingAccount struct{}

func (InstructionError_MissingAccount) Error() string {
	return "An account required by the instruction is missing"
}
func (InstructionError_MissingAccount) xx_isInstructionError() {}

// Cross-program invocation reentrancy not allowed for this instruction
type InstructionError_ReentrancyNotAllowed struct{}

func (InstructionError_ReentrancyNotAllowed) Error() string {
	return "Cross-program invocation reentrancy not allowed for this instruction"
}
func (InstructionError_ReentrancyNotAllowed) xx_isInstructionError() {}

// Length of the seed is too long for address generation
type InstructionError_MaxSeedLengthExceeded struct{}

func (InstructionError_MaxSeedLengthExceeded) Error() string {
	return "Length of the seed is too long for address generation"
}
func (InstructionError_MaxSeedLengthExceeded) xx_isInstructionError() {}

// Provided seeds do not result in a valid address
type InstructionError_InvalidSeeds struct{}

func (InstructionError_InvalidSeeds) Error() string {
	return "Provided seeds do not result in a valid address"
}
func (InstructionError_InvalidSeeds) xx_isInstructionError() {}

// Failed to reallocate account data of this length
type InstructionError_InvalidRealloc struct{}

func (InstructionError_InvalidRealloc) Error() string          { return "Failed to reallocate account data" }
func (InstructionError_InvalidRealloc) xx_isInstructionError() {}

// Computational budget exceeded
type InstructionError_ComputationalBudgetExceeded struct{}

func (InstructionError_ComputationalBudgetExceeded) Error() string {
	return "Computational budget exceeded"
}
func (InstructionError_ComputationalBudgetExceeded) xx_isInstructionError() {}

// Cross-program invocation with unauthorized signer or writable account
type InstructionError_PrivilegeEscalation struct{}

func (InstructionError_PrivilegeEscalation) Error() string {
	return "Cross-program invocation with unauthorized signer or writable account"
}
func (InstructionError_PrivilegeEscalation) xx_isInstructionError() {}

// Failed to create program execution environment
type InstructionError_ProgramEnvironmentSetupFailure struct{}

func (InstructionError_ProgramEnvironmentSetupFailure) Error() string {
	return "Failed to create program execution environment"
}
func (InstructionError_ProgramEnvironmentSetupFailure) xx_isInstructionError() {}

// Program failed to complete
type InstructionError_ProgramFailedToComplete struct{}

func (InstructionError_ProgramFailedToComplete) Error() string          { return "Program failed to complete" }
func (InstructionError_ProgramFailedToComplete) xx_isInstructionError() {}

// Program failed to compile
type InstructionError_ProgramFailedToCompile struct{}

func (InstructionError_ProgramFailedToCompile) Error() string          { return "Program failed to compile" }
func (InstructionError_ProgramFailedToCompile) xx_isInstructionError() {}

// Account is immutable
type InstructionError_Immutable struct{}

func (InstructionError_Immutable) Error() string          { return "Account is immutable" }
func (InstructionError_Immutable) xx_isInstructionError() {}

// Incorrect authority provided
type InstructionError_IncorrectAuthority struct{}

func (InstructionError_IncorrectAuthority) Error() string          { return "Incorrect authority provided" }
func (InstructionError_IncorrectAuthority) xx_isInstructionError() {}

// An account does not have enough lamports to be rent-exempt
type InstructionError_AccountNotRentExempt struct{}

func (InstructionError_AccountNotRentExempt) Error() string {
	return "An account does not have enough lamports to be rent-exempt"
}
func (InstructionError_AccountNotRentExempt) xx_isInstructionError() {}

// Invalid account owner
type InstructionError_InvalidAccountOwner struct{}

func (InstructionError_InvalidAccountOwner) Error() string          { return "Invalid account owner" }
func (InstructionError_InvalidAccountOwner) xx_isInstructionError() {}

// Program arithmetic overflowed
type InstructionError_ArithmeticOverflow struct{}

func (InstructionError_ArithmeticOverflow) Error() string          { return "Program arithmetic overflowed" }
func (InstructionError_ArithmeticOverflow) xx_isInstructionError() {}

// Unsupported sysvar
type InstructionError_UnsupportedSysvar struct{}

func (InstructionError_UnsupportedSysvar) Error() string          { return "Unsupported sysvar" }
func (InstructionError_UnsupportedSysvar) xx_isInstructionError() {}

// Illegal account owner
type InstructionError_IllegalOwner struct{}

func (InstructionError_IllegalOwner) Error() string          { return "Provided owner is not allowed" }
func (InstructionError_IllegalOwner) xx_isInstructionError() {}
