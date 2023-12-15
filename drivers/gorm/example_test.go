//go:build go1.16
// +build go1.16

package gorm_test

import (
	"context"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	trmgorm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	trm "github.com/avito-tech/go-transaction-manager/v2"
	"github.com/avito-tech/go-transaction-manager/v2/manager"
	"github.com/avito-tech/go-transaction-manager/v2/settings"
)

// Example demonstrates the implementation of the Repository pattern by trm.Manager.
func Example() {
	db, err := gorm.Open(sqlite.Open("file:test.drivers?mode=memory"))
	checkErr(err)

	// Migrate the schema
	checkErr(db.AutoMigrate(&userRow{}))

	r := newRepo(db, trmgorm.DefaultCtxGetter)

	u := &user{
		Username: "username",
	}

	ctx := context.Background()
	trManager := manager.Must(
		trmgorm.NewDefaultFactory(db),
		manager.WithSettings(trmgorm.MustSettings(
			settings.Must(
				settings.WithPropagation(trm.PropagationNested))),
		),
	)

	err = trManager.Do(ctx, func(ctx context.Context) error {
		if err := r.Save(ctx, u); err != nil {
			return err
		}

		return trManager.Do(ctx, func(ctx context.Context) error {
			u.Username = "new_username"

			return r.Save(ctx, u)
		})
	})
	checkErr(err)

	userFromDB, err := r.GetByID(ctx, u.ID)
	checkErr(err)

	fmt.Println(userFromDB)

	// Output: &{1 new_username}
}

type repo struct {
	db     *gorm.DB
	getter *trmgorm.CtxGetter
}

func newRepo(db *gorm.DB, c *trmgorm.CtxGetter) *repo {
	return &repo{
		db:     db,
		getter: c,
	}
}

type user struct {
	ID       int64
	Username string
}

type userRow struct {
	ID       int64 `gorm:"primarykey"`
	Username string
}

func (r *repo) GetByID(ctx context.Context, id int64) (*user, error) {
	var row userRow
	db := r.getter.DefaultTrOrDB(ctx, r.db).
		WithContext(ctx).Model(userRow{ID: id}).First(&row)

	if db.Error != nil {
		return nil, db.Error
	}

	return r.toModel(row), nil
}

func (r *repo) Save(ctx context.Context, u *user) error {
	isNew := u.ID == 0

	db := r.getter.DefaultTrOrDB(ctx, r.db).WithContext(ctx)

	row := r.toRow(u)
	if isNew {
		db = db.Create(&row)

		u.ID = row.ID
	} else {
		db = db.Save(&row)
	}

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (r *repo) toRow(model *user) userRow {
	return userRow{
		ID:       model.ID,
		Username: model.Username,
	}
}

func (r *repo) toModel(row userRow) *user {
	return &user{
		ID:       row.ID,
		Username: row.Username,
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
