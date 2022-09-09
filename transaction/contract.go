// Package transaction is an interface to create a transactional usecase
// in the Application layer.
package transaction

import (
	"context"
	"errors"
	"fmt"
)

// ErrTransaction is an error while working with a transaction.
var (
	ErrTransaction = errors.New("transaction")
	ErrBegin       = errTransaction("begin")
	ErrCommit      = errTransaction("close")
	ErrRollback    = errTransaction("rollback")
)

func errTransaction(msg string) error {
	return fmt.Errorf("%w: %s", ErrTransaction, msg)
}

// Manager manages a transaction from Begin to Commit or Rollback.
type Manager interface {
	// Do processes a transaction inside a closure.
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// Factory is used in Manager to creates Transaction.
// TODO probably Settings is required as the func argument.
type Factory func() (Transaction, error)

// Transaction wraps different transaction implementations.
type Transaction interface {
	// Transaction returns the real transaction sql.Tx, sqlx.Tx or another.
	Transaction() interface{}
	// Commit calls close for a database.
	Commit() error
	// Rollback calls close for a database.
	Rollback() error
	// IsActive returns true if the transaction started but not committed or rolled back.
	IsActive() bool
}

// Settings is configuration for Manager.
// TODO probably needs to separate Transaction and Manager settings.
type Settings interface {
	CtxKey() CtxKey
	IsReadOnly() bool
	Propagation() Propagation
}

// Propagation is a type for transaction propagation rules.
type Propagation int8

// TODO fix description and implement
// now is copy of
//
//nolint:lll https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/transaction/annotation/Propagation.html
const (
	// PropagationRequired supports a current transaction, create a new one if none exists.
	PropagationRequired Propagation = iota
	// PropagationNested executes within a nested transaction
	// if a current transaction exists, behave like REQUIRED otherwise.
	PropagationNested
	// PropagationsMandatory supports a current transaction.
	PropagationsMandatory
	// PropagationNever executes non-transactionally, throw an exception if a transaction exists.
	PropagationNever
	// PropagationNotSupported executes non-transactionally, suspend the current transaction if one exists..
	PropagationNotSupported
	// PropagationRequiresNew creates a new transaction, and suspend the current transaction if one exists.
	PropagationRequiresNew
	// PropagationSupports supports a current transaction, execute non-transactionally if none exists.
	PropagationSupports
)
