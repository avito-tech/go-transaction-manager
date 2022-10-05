// Package transaction is an interface to create a transactional usecase
// in the Application layer.
package transaction

import (
	"errors"
	"fmt"
)

var (
	// ErrTransaction is an error while working with a transaction.
	ErrTransaction = errors.New("transaction")
	// ErrCommit occurs when commit finished with an error.
	ErrCommit = errTransaction("close")
	// ErrRollback occurs when rollback finished with an error.
	ErrRollback = errTransaction("rollback")
)

func errNested(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}

func errTransaction(msg string) error {
	return errNested(ErrTransaction, msg)
}

// TrFactory is used in Manager to creates Transaction.
type TrFactory func() (Transaction, error)

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

var (
	// ErrPropagation occurs because of Propagation setting.
	ErrPropagation = errTransaction("propagation")
	// TODO need implement.
	ErrPropagationMandatory = errNested(ErrPropagation, "mandatory")
	// TODO need implement.
	ErrPropagationNever = errNested(ErrPropagation, "never")
)

// Propagation is a type for transaction propagation rules.
type Propagation int8

// TODO fix description and implement, there is not for NoSQL
// now is copy of
//
// https://docs.spring.io/spring-framework/docs/current/javadoc-api/org/springframework/transaction/annotation/Propagation.html //nolint:lll
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
	// PropagationNotSupported executes non-transactionally, suspend the current transaction if one exists.
	PropagationNotSupported
	// PropagationRequiresNew creates a new transaction, and suspend the current transaction if one exists.
	PropagationRequiresNew
	// PropagationSupports supports a current transaction, execute non-transactionally if none exists.
	PropagationSupports
)
