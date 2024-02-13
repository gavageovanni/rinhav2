package config2

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db *pgxpool.Pool
)

func Init() error {
	var err error
	db, err = getPostgresConnection()
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *pgxpool.Pool {
	return db
}
