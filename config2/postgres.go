package config2

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var (
	host  = "localhost"
	port  = 5432
	user  = "admin"
	pass  = "123"
	dbase = "rinha"
)

func getPostgresConnection() (*pgxpool.Pool, error) {
	psqlConnStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbase)
	poolConfig, err := pgxpool.ParseConfig(psqlConnStr)

	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return db, nil
}
