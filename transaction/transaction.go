// Package transaction is an interface to create a transactional usecase
// in the Application layer.
package transaction

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrTransaction is an error while working with a transaction.
	ErrTransaction = errors.New("transaction")
	// ErrAlreadyClosed occurs if the transaction was closed outside the Manager.
	ErrAlreadyClosed = errTransaction("already closed")

	// ErrBegin occurs when a transaction started with an error.
	ErrBegin = errTransaction("begin")
	// ErrCommit occurs when commit finished with an error.
	ErrCommit = errTransaction("commit")
	// ErrRollback occurs when rollback finished with an error.
	ErrRollback = errTransaction("rollback")

	// ErrSPBegin occurs when a savepoint started with an error.
	ErrSPBegin = errNested(ErrBegin, "savepoint")
	// ErrSPCommit occurs when release savepoint finished with an error.
	ErrSPCommit = errNested(ErrCommit, "savepoint")
	// ErrSPRollback occurs when rollback savepoint finished with an error.
	ErrSPRollback = errNested(ErrRollback, "savepoint")
)

func errNested(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

func errTransaction(msg string) error {
	return errNested(ErrTransaction, msg)
}

// TrFactory is used in Manager to creates Transaction.
type TrFactory func(ctx context.Context, s Settings) (context.Context, Transaction, error)

// SPFactory creates save points for Transaction.
type SPFactory interface {
	SavePoint(ctx context.Context, s Settings) (context.Context, Transaction, error)
}

// Transaction wraps different transaction implementations.
type Transaction interface {
	// Transaction returns the real transaction sql.Tx, sqlx.Tx or another.
	Transaction() interface{}
	// Commit the transaction.Transaction.
	Commit(context.Context) error
	// Rollback the transaction.Transaction.
	Rollback(context.Context) error
	// IsActive returns true if the transaction started but not committed or rolled back.
	IsActive() bool
}

// transactionWithSP is used for tests.
//
//nolint:unused
type transactionWithSP interface {
	Transaction
	SPFactory
}

var (
	// ErrPropagation occurs because of Propagation setting.
	ErrPropagation = errTransaction("propagation")
	// ErrPropagationMandatory occurs when the transaction doesn't exist.
	ErrPropagationMandatory = errNested(ErrPropagation, "mandatory")
	// ErrPropagationNever occurs when the transaction already exists.
	ErrPropagationNever = errNested(ErrPropagation, "never")
)

// Propagation is a type for transaction propagation rules.
type Propagation int8

const (
	// PropagationRequired supports a current transaction, create a new one if none exists. This is default setting.
	PropagationRequired Propagation = iota
	// PropagationNested executes within a nested transaction
	// if a current transaction exists, create a new one if none exists.
	PropagationNested
	// PropagationsMandatory supports a current transaction, throws an exception if none exists.
	PropagationsMandatory
	// PropagationNever executes non-transactionally, throws an exception if a transaction exists.
	PropagationNever
	// PropagationNotSupported executes non-transactionally, suspends the current transaction if one exists.
	PropagationNotSupported
	// PropagationRequiresNew creates a new transaction, suspends the current transaction if one exists.
	PropagationRequiresNew
	// PropagationSupports supports a current transaction, execute non-transactionally if none exists.
	PropagationSupports
)
