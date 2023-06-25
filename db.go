package sua

import "database/sql"

func NewDB(sqldb *sql.DB) *DB {
	db := &DB{DB: sqldb}
	db.stmt = stmt{DB: db}
	return db
}

type DB struct {
	*sql.DB
	stmt
}
