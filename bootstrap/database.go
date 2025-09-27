package bootstrap

import (
	"database/sql"
	"sync"

	"github.com/aarondl/sqlboiler/v4/boil"
	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

func SetUpDb(dsn string) (*sql.DB, error) {

	var err error

	once.Do(func() {

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			return
		}

		// if err := db.Ping(); err != nil {
		// 	return nil, err
		// }

		boil.SetDB(db)
	})

	return db, err
}
