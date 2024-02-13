package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"os"
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

var (
	host  = os.Getenv("POSTGRES_HOST")
	port  = 5432
	user  = "rinhav2"
	pass  = "rinhav2"
	dbase = "rinhav2"
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
