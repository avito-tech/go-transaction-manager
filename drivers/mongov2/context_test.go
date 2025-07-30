//go:build go1.21

package mongov2

import (
	"context"
	"testing"

	"github.com/avito-tech/go-transaction-manager/drivers/mongov2/v2/internal/mtest"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/avito-tech/go-transaction-manager/trm/v2/settings"
	"github.com/stretchr/testify/require"
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
