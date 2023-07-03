package token

type Error interface {
	error
	xx_isTokenError()
}

func CustomErrorResolver(code int) (error, bool) {
	switch code {
	case 0:
		return Error_NotRentExempt{}, true
	case 1:
		return Error_InsufficientFunds{}, true
	case 2:
		return Error_InvalidMint{}, true
	case 3:
		return Error_MintMismatch{}, true
	case 4:
		return Error_OwnerMismatch{}, true
	case 5:
		return Error_FixedSupply{}, true
	case 6:
		return Error_AlreadyInUse{}, true
	case 7:
		return Error_InvalidNumberOfProvidedSigners{}, true
	case 8:
		return Error_InvalidNumberOfRequiredSigners{}, true
	case 9:
		return Error_UninitializedState{}, true
	case 10:
		return Error_NativeNotSupported{}, true
	case 11:
		return Error_NonNativeHasBalance{}, true
	case 12:
		return Error_InvalidInstruction{}, true
	case 13:
		return Error_InvalidState{}, true
	case 14:
		return Error_Overflow{}, true
	case 15:
		return Error_AuthorityTypeNotSupported{}, true
	case 16:
		return Error_MintCannotFreeze{}, true
	case 17:
		return Error_AccountFrozen{}, true
	case 18:
		return Error_MintDecimalsMismatch{}, true
	case 19:
		return Error_NonNativeNotSupported{}, true
	default:
		return nil, false
	}
}

// Defined [here](https://github.com/solana-labs/solana-program-library/blob/9b3b5d8841484f19f927b1f49aa37bbd4da2a1ca/token/program/src/error.rs#L13).

// Lamport balance below rent-exempt threshold.
type Error_NotRentExempt struct{}

func (Error_NotRentExempt) Error() string {
	return "Lamport balance below rent-exempt threshold"
}

func (Error_NotRentExempt) xx_isTokenError() {}

// Insufficient funds for the operation requested.
type Error_InsufficientFunds struct{}

func (Error_InsufficientFunds) Error() string {
	return "Insufficient funds"
}

func (Error_InsufficientFunds) xx_isTokenError() {}

// Invalid Mint.
type Error_InvalidMint struct{}

func (Error_InvalidMint) Error() string {
	return "Invalid Mint"
}

func (Error_InvalidMint) xx_isTokenError() {}

// Account not associated with this Mint.
type Error_MintMismatch struct{}

func (Error_MintMismatch) Error() string {
	return "Account not associated with this Mint"
}

func (Error_MintMismatch) xx_isTokenError() {}

// Owner does not match.
type Error_OwnerMismatch struct{}

func (Error_OwnerMismatch) Error() string {
	return "Owner does not match"
}

func (Error_OwnerMismatch) xx_isTokenError() {}

// This token's supply is fixed and new tokens cannot be minted.
type Error_FixedSupply struct{}

func (Error_FixedSupply) Error() string {
	return "Fixed supply"
}

func (Error_FixedSupply) xx_isTokenError() {}

// The account cannot be initialized because it is already being used.
type Error_AlreadyInUse struct{}

func (Error_AlreadyInUse) Error() string {
	return "Already in use"
}

func (Error_AlreadyInUse) xx_isTokenError() {}

// Invalid number of provided signers.
type Error_InvalidNumberOfProvidedSigners struct{}

func (Error_InvalidNumberOfProvidedSigners) Error() string {
	return "Invalid number of provided signers"
}

func (Error_InvalidNumberOfProvidedSigners) xx_isTokenError() {}

// Invalid number of required signers.
type Error_InvalidNumberOfRequiredSigners struct{}

func (Error_InvalidNumberOfRequiredSigners) Error() string {
	return "Invalid number of required signers"
}

func (Error_InvalidNumberOfRequiredSigners) xx_isTokenError() {}

// State is uninitialized.
type Error_UninitializedState struct{}

func (Error_UninitializedState) Error() string {
	return "State is uninitialized"
}

func (Error_UninitializedState) xx_isTokenError() {}

// Instruction does not support native tokens.
type Error_NativeNotSupported struct{}

func (Error_NativeNotSupported) Error() string {
	return "Instruction does not support native tokens"
}

func (Error_NativeNotSupported) xx_isTokenError() {}

// Non-native account can only be closed if its balance is zero.
type Error_NonNativeHasBalance struct{}

func (Error_NonNativeHasBalance) Error() string {
	return "Non-native account can only be closed if its balance is zero"
}

func (Error_NonNativeHasBalance) xx_isTokenError() {}

// Invalid instruction.
type Error_InvalidInstruction struct{}

func (Error_InvalidInstruction) Error() string {
	return "Invalid instruction"
}

func (Error_InvalidInstruction) xx_isTokenError() {}

// State is invalid for requested operation.
type Error_InvalidState struct{}

func (Error_InvalidState) Error() string {
	return "State is invalid for requested operation"
}

func (Error_InvalidState) xx_isTokenError() {}

// Operation overflowed.
type Error_Overflow struct{}

func (Error_Overflow) Error() string {
	return "Operation overflowed"
}

func (Error_Overflow) xx_isTokenError() {}

// Account does not support specified authority type.
type Error_AuthorityTypeNotSupported struct{}

func (Error_AuthorityTypeNotSupported) Error() string {
	return "Account does not support specified authority type"
}

func (Error_AuthorityTypeNotSupported) xx_isTokenError() {}

// This token mint cannot freeze accounts.
type Error_MintCannotFreeze struct{}

func (Error_MintCannotFreeze) Error() string {
	return "This token mint cannot freeze accounts"
}

func (Error_MintCannotFreeze) xx_isTokenError() {}

// Account is frozen; all account operations will fail.
type Error_AccountFrozen struct{}

func (Error_AccountFrozen) Error() string {
	return "Account is frozen"
}

func (Error_AccountFrozen) xx_isTokenError() {}

// Mint decimals mismatch between the client and mint.
type Error_MintDecimalsMismatch struct{}

func (Error_MintDecimalsMismatch) Error() string {
	return "The provided decimals value different from the Mint decimals"
}

func (Error_MintDecimalsMismatch) xx_isTokenError() {}

// Instruction does not support non-native tokens.
type Error_NonNativeNotSupported struct{}

func (Error_NonNativeNotSupported) Error() string {
	return "Instruction does not support non-native tokens"
}

func (Error_NonNativeNotSupported) xx_isTokenError() {}
