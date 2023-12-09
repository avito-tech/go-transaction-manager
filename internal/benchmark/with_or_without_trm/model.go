package with_or_without_trm

type user struct {
	ID       int64  `dbs:"user_id"`
	Username string `dbs:"username"`
}
