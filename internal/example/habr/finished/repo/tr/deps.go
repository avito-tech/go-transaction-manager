package tr

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Tr interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	sqlx.Ext
	sqlx.ExtContext
}
