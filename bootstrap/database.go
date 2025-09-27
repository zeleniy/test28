package bootstrap

import (
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	_ "github.com/lib/pq"
)

func SetUpDb(dsn string) (*sql.DB, error) {

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	// if err := db.Ping(); err != nil {
	// 	return nil, err
	// }

	boil.SetDB(db)

	return db, nil
}
