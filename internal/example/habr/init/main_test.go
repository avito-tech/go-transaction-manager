package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"init/domain"
	"init/queue"
	"init/repo"
	"init/usecase/fast_purchase"
	purchaseUC "init/usecase/purchase"
	registerUC "init/usecase/register"
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
	userRepo := repo.NewUserRepo(db)
	orderRepo := repo.NewOrdrRepo(db)

	queueRegistered := queue.Queue[domain.Registered]{}
	queuePurchased := queue.Queue[domain.Purchased]{}

	// usecase
	register := registerUC.New(userRepo, db, queueRegistered)
	purchase := purchaseUC.New(orderRepo, db, queuePurchased)
	fastPurchase := fast_purchase.New(db, register, purchase)

	out, err := fastPurchase.Handle(fast_purchase.In{
		Register: fast_purchase.RegisterIn{
			Username: "username",
			Password: "password",
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
