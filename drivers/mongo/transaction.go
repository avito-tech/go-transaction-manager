// Package mongo is an implementation of trm.Transaction interface by Transaction for mongo.Client.
package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers"
)

// Transaction is trm.Transaction for mongo.Client.
type Transaction struct {
	session  mongo.Session
	isClosed *drivers.IsClosed
}

// NewTransaction creates trm.Transaction for mongo.Client.
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

	tr := &Transaction{session: s, isClosed: drivers.NewIsClosed()}

	go tr.awaitDone(ctx)

	return mongo.NewSessionContext(ctx, tr.session), tr, nil
}

func (t *Transaction) awaitDone(ctx context.Context) {
	if ctx.Done() == nil {
		return
	}

	select {
	case <-ctx.Done():
		t.isClosed.Close()
	case <-t.isClosed.Closed():
	}
}

// Transaction returns the real transaction mongo.Session.
func (t *Transaction) Transaction() interface{} {
	return t.session
}

// Commit the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	defer t.isClosed.Close()

	defer t.session.EndSession(ctx)

	return t.session.CommitTransaction(ctx)
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	defer t.isClosed.Close()

	defer t.session.EndSession(ctx)

	return t.session.AbortTransaction(ctx)
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return t.isClosed.IsActive()
}

// Closed returns a channel that's closed when transaction committed or rolled back.
func (t *Transaction) Closed() <-chan struct{} {
	return t.isClosed.Closed()
}
