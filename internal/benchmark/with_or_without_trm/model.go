package with_or_without_trm

type user struct {
	ID       int64  `db:"user_id"`
	Username string `db:"username"`
}
