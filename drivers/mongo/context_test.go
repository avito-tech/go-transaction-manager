//go:build go1.14
// +build go1.14

package mongo

import (
	"context"
	"testing"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mt := mtest.New(
		t,
		mtest.NewOptions().ClientType(mtest.Mock),
	)

	mt.Run("all", func(mt *mtest.T) {
		mt.Parallel()

		m := manager.Must(
			NewDefaultFactory(mt.Client),
		)

		err := m.Do(ctx, func(ctx context.Context) error {
			tr := DefaultCtxGetter.TrOrDB(ctx, settings.DefaultCtxKey, nil)
			require.NotNil(t, tr)

			tr = DefaultCtxGetter.DefaultTrOrDB(ctx, nil)
			require.NotNil(t, tr)

			tr = DefaultCtxGetter.TrOrDB(ctx, "invalid ley", nil)
			require.Nil(t, tr)

			err := m.Do(ctx, func(ctx context.Context) error {
				tr = DefaultCtxGetter.DefaultTrOrDB(ctx, nil)
				require.NotNil(t, tr)

				tr = DefaultCtxGetter.TrOrDB(ctx, settings.DefaultCtxKey, nil)
				require.NotNil(t, tr)

				tr = DefaultCtxGetter.TrOrDB(ctx, "invalid ley", nil)
				require.Nil(t, tr)

				return nil
			})

			return err
		})

		require.NoError(t, err)
	})
}
