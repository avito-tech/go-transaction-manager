package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	trmsqlx "github.com/avito-tech/go-transaction-manager/sqlx"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"finish/domain"
	"finish/queue"
	"finish/repo"
	"finish/usecase/fast_purchase"
	purchaseUC "finish/usecase/purchase"
	registerUC "finish/usecase/register"
)

// language=SQLite
const sqlStmt = `
DROP TABLE IF EXISTS "user";
CREATE TABLE IF NOT EXISTS "user"
(
    id           INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    username     TEXT    NOT NULL,
    password     TEXT    NOT NULL,
    notification TEXT
);
INSERT INTO "user" (username, password, notification) VALUES ('user1', 'password1', 'email1');

DROP TABLE IF EXISTS "order";
CREATE TABLE IF NOT EXISTS "order"
(
    id         INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    user_id    INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity   INTEGER NOT NULL
);
`

func Example() {
	db, err := sqlx.Open("sqlite3", "file:test.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	// Repo
	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	userRepo := repo.NewUserRepo(db, trmsqlx.DefaultCtxGetter)
	orderRepo := repo.NewOrderRepo(db, trmsqlx.DefaultCtxGetter)

	queueRegistered := queue.Queue[domain.Registered]{}
	queuePurchased := queue.Queue[domain.Purchased]{}

	// usecase
	register := registerUC.New(userRepo, trManager, queueRegistered)
	purchase := purchaseUC.New(orderRepo, trManager, queuePurchased)
	fastPurchase := fast_purchase.New(trManager, register, purchase)

	ctx := context.Background()

	out, err := fastPurchase.Handle(ctx, fast_purchase.In{
		Register: fast_purchase.RegisterIn{
			Username: "username",
		},
		Purchase: fast_purchase.PurchaseIn{
			ProductID: 1,
			Quantity:  10,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	res, _ := json.MarshalIndent(out, "", "\t")
	fmt.Println(string(res))

	// Output: {
	//	"User": {
	//		"ID": 2,
	//		"Username": "username",
	//		"Password": "",
	//		"Notification": {
	//			"Email": false,
	//			"SMS": true
	//		}
	//	},
	//	"Order": {
	//		"ID": 1,
	//		"ProductID": 1,
	//		"UserID": 2,
	//		"Quantity": 10
	//	}
	//}
}
