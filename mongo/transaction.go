// Package mongo is an implementation of transaction.Transaction interface by Transaction for mongo.Client.
package mongo

import (
	"context"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Transaction is transaction.Transaction for mongo.Client.
type Transaction struct {
	session  mongo.Session
	isActive int64
}

// NewTransaction creates transaction.Transaction for mongo.Client.
func NewTransaction(
	ctx context.Context,
	sessionOptions *options.SessionOptions,
	trOpts *options.TransactionOptions,
	client client,
) (context.Context, *Transaction, error) {
	s, err := client.StartSession(sessionOptions)
	if err != nil {
		return ctx, nil, err
	}

	if err = s.StartTransaction(trOpts); err != nil {
		defer s.EndSession(ctx)

		return ctx, nil, err
	}

	tr := &Transaction{session: s, isActive: 1}

	go tr.awaitDone(ctx)

	return mongo.NewSessionContext(ctx, tr.session), tr, nil
}

func (t *Transaction) awaitDone(ctx context.Context) {
	if ctx.Done() == nil {
		return
	}

	<-ctx.Done()

	t.deactivate()
}

// Transaction returns the real transaction mongo.Session.
func (t *Transaction) Transaction() interface{} {
	return t.session
}

// Commit the transaction.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	t.deactivate()

	defer t.session.EndSession(ctx)

	if err := t.session.CommitTransaction(ctx); err != nil {
		return err
	}

	return nil
}

// Rollback the transaction.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	t.deactivate()

	defer t.session.EndSession(ctx)

	if err := t.session.AbortTransaction(ctx); err != nil {
		return err
	}

	return nil
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return atomic.LoadInt64(&t.isActive) == 1
}

func (t *Transaction) deactivate() {
	atomic.SwapInt64(&t.isActive, 0)
}
