package solana

import (
	jsn "encoding/json"
	"fmt"
)

type TransactionError interface {
	error
	RawError() interface{}
	xx_isTransactionError()
}

type transactionError struct {
	cause error
	raw   interface{}
}

var _ TransactionError = &transactionError{}

func (v *transactionError) Error() string {
	return v.cause.Error()
}

func (v *transactionError) Unwrap() error {
	return v.cause
}

func (v *transactionError) Cause() error {
	return v.cause
}

func (v *transactionError) RawError() interface{} {
	return v.raw
}

func (transactionError) xx_isTransactionError() {}

func ParseTransactionError(tx *Transaction, raw interface{}) (TransactionError, bool) {
	switch t := raw.(type) {
	case string:
		err, ok := parseTransactionErrorString(t)
		if !ok {
			return nil, false
		}
		return &transactionError{
			cause: err,
			raw:   raw,
		}, true
	case map[string]interface{}:
		err, ok := parseTransactionErrorObject(tx, t)
		if !ok {
			return nil, false
		}
		return &transactionError{
			cause: err,
			raw:   raw,
		}, true
	default:
		return nil, false
	}
}

func parseTransactionErrorObject(tx *Transaction, err map[string]interface{}) (error, bool) {
	fields, ok := err["InstructionError"].([]interface{})
	if !ok {
		return nil, false
	}
	if len(fields) != 2 {
		return nil, false
	}
	index, ok := asFloat64(fields[0])
	if !ok {
		return nil, false
	}
	var progID *PublicKey
	if tx != nil {
		in := tx.Message.Instructions[int(index)]
		prog, rErr := tx.ResolveProgramIDIndex(in.ProgramIDIndex)
		if rErr != nil {
			return nil, false //nolint: nilerr
		}
		progID = &prog
	}
	cause, ok := ParseInstructionError(fields[1], progID)
	if !ok {
		return nil, false
	}
	return &TransactionError_InstructionError{
		Index: int32(index),
		Cause: cause,
	}, true
}

func asFloat64(v interface{}) (float64, bool) {
	index, ok := v.(float64)
	if ok {
		return index, true
	}
	s, ok := v.(jsn.Number)
	if !ok {
		return 0, false
	}
	index, err := s.Float64()
	if err != nil {
		return 0, false
	}
	return index, true
}

func parseTransactionErrorString(err string) (error, bool) {
	switch err {
	case "AccountInUse":
		return TransactionError_AccountInUse{}, true
	case "AccountLoadedTwice":
		return TransactionError_AccountLoadedTwice{}, true
	case "AccountNotFound":
		return TransactionError_AccountNotFound{}, true
	case "ProgramAccountNotFound":
		return TransactionError_ProgramAccountNotFound{}, true
	case "InsufficientFundsForFee":
		return TransactionError_InsufficientFundsForFee{}, true
	case "InvalidAccountForFee":
		return TransactionError_InvalidAccountForFee{}, true
	case "AlreadyProcessed":
		return TransactionError_AlreadyProcessed{}, true
	case "BlockhashNotFound":
		return TransactionError_BlockhashNotFound{}, true
	case "CallChainTooDeep":
		return TransactionError_CallChainTooDeep{}, true
	case "MissingSignatureForFee":
		return TransactionError_MissingSignatureForFee{}, true
	case "InvalidAccountIndex":
		return TransactionError_InvalidAccountIndex{}, true
	case "SignatureFailure":
		return TransactionError_SignatureFailure{}, true
	case "InvalidProgramForExecution":
		return TransactionError_InvalidProgramForExecution{}, true
	case "SanitizeFailure":
		return TransactionError_SanitizeFailure{}, true
	case "ClusterMaintenance":
		return TransactionError_ClusterMaintenance{}, true
	case "AccountBorrowOutstanding":
		return TransactionError_AccountBorrowOutstanding{}, true
	default:
		return TransactionError_Undefined(err), true
	}
}

type TransactionError_Undefined string

func (v TransactionError_Undefined) Error() string {
	return string(v)
}
func (TransactionError_Undefined) xx_isTransactionError() {}

// Defined [here](https://github.com/solana-labs/solana/blob/c0c60386544ec9a9ec7119229f37386d9f070523/sdk/src/transaction/error.rs#L13).

// An account is already being processed in another transaction in a way that
// does not support parallelism
type TransactionError_AccountInUse struct{}

func (TransactionError_AccountInUse) Error() string          { return "Account in use" }
func (TransactionError_AccountInUse) xx_isTransactionError() {}

// A `Pubkey` appears twice in the transaction's `account_keys`.  Instructions
// can reference `Pubkey`s more than once but the message must contain a list
// with no duplicate keys
type TransactionError_AccountLoadedTwice struct{}

func (TransactionError_AccountLoadedTwice) Error() string          { return "Account loaded twice" }
func (TransactionError_AccountLoadedTwice) xx_isTransactionError() {}

// Attempt to debit an account but found no record of a prior credit.
type TransactionError_AccountNotFound struct{}

func (TransactionError_AccountNotFound) Error() string {
	return "Attempt to debit an account but found no record of a prior credit."
}
func (TransactionError_AccountNotFound) xx_isTransactionError() {}

// Attempt to load a program that does not exist
type TransactionError_ProgramAccountNotFound struct{}

func (TransactionError_ProgramAccountNotFound) Error() string {
	return "Attempt to load a program that does not exist"
}
func (TransactionError_ProgramAccountNotFound) xx_isTransactionError() {}

// The from `Pubkey` does not have sufficient balance to pay the fee to
// schedule the transaction
type TransactionError_InsufficientFundsForFee struct{}

func (TransactionError_InsufficientFundsForFee) Error() string          { return "Insufficient funds for fee" }
func (TransactionError_InsufficientFundsForFee) xx_isTransactionError() {}

// This account may not be used to pay transaction fees
type TransactionError_InvalidAccountForFee struct{}

func (TransactionError_InvalidAccountForFee) Error() string {
	return "This account may not be used to pay transaction fees"
}
func (TransactionError_InvalidAccountForFee) xx_isTransactionError() {}

// The bank has seen this transaction before. This can occur under normal
// operation when a UDP packet is duplicated, as a user error from a client not
// updating its `recent_blockhash`, or as a double-spend attack.
type TransactionError_AlreadyProcessed struct{}

func (TransactionError_AlreadyProcessed) Error() string {
	return "This transaction has already been processed"
}
func (TransactionError_AlreadyProcessed) xx_isTransactionError() {}

// The bank has not seen the given `recent_blockhash` or the transaction is too
// old and the `recent_blockhash` has been discarded.
type TransactionError_BlockhashNotFound struct{}

func (TransactionError_BlockhashNotFound) Error() string          { return "Blockhash not found" }
func (TransactionError_BlockhashNotFound) xx_isTransactionError() {}

// An error occurred while processing an instruction.
type TransactionError_InstructionError struct {
	Index int32
	Cause InstructionError
}

func (v *TransactionError_InstructionError) Error() string {
	return fmt.Sprintf("Error processing instruction %d: %s", v.Index, v.Cause.Error())
}
func (v *TransactionError_InstructionError) Unwrap() error        { return v.Cause }
func (*TransactionError_InstructionError) xx_isTransactionError() {}

// Loader call chain is too deep
type TransactionError_CallChainTooDeep struct{}

func (TransactionError_CallChainTooDeep) Error() string          { return "Loader call chain is too deep" }
func (TransactionError_CallChainTooDeep) xx_isTransactionError() {}

// Transaction requires a fee but has no signature present
type TransactionError_MissingSignatureForFee struct{}

func (TransactionError_MissingSignatureForFee) Error() string {
	return "Transaction requires a fee but has no signature present"
}
func (TransactionError_MissingSignatureForFee) xx_isTransactionError() {}

// Transaction contains an invalid account reference
type TransactionError_InvalidAccountIndex struct{}

func (TransactionError_InvalidAccountIndex) Error() string {
	return "Transaction contains an invalid account reference"
}
func (TransactionError_InvalidAccountIndex) xx_isTransactionError() {}

// Transaction did not pass signature verification
type TransactionError_SignatureFailure struct{}

func (TransactionError_SignatureFailure) Error() string {
	return "Transaction did not pass signature verification"
}
func (TransactionError_SignatureFailure) xx_isTransactionError() {}

// This program may not be used for executing instructions
type TransactionError_InvalidProgramForExecution struct{}

func (TransactionError_InvalidProgramForExecution) Error() string {
	return "This program may not be used for executing instructions"
}
func (TransactionError_InvalidProgramForExecution) xx_isTransactionError() {}

// Transaction failed to sanitize accounts offsets correctly implies that
// account locks are not taken for this TX, and should not be unlocked.
type TransactionError_SanitizeFailure struct{}

func (TransactionError_SanitizeFailure) Error() string {
	return "Transaction failed to sanitize accounts offsets correctly"
}
func (TransactionError_SanitizeFailure) xx_isTransactionError() {}

type TransactionError_ClusterMaintenance struct{}

func (TransactionError_ClusterMaintenance) Error() string {
	return "Transactions are currently disabled due to cluster maintenance"
}
func (TransactionError_ClusterMaintenance) xx_isTransactionError() {}

// Transaction processing left an account with an outstanding borrowed reference
type TransactionError_AccountBorrowOutstanding struct{}

func (TransactionError_AccountBorrowOutstanding) Error() string {
	return "Transaction processing left an account with an outstanding borrowed reference"
}
func (TransactionError_AccountBorrowOutstanding) xx_isTransactionError() {}
