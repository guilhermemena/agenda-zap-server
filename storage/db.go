package storage

import (
	"context"
	"log"

	"github.com/guilhermemena/agenda-zap-server/cmd/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func connectToDB() *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), configs.Envs.DBConnection)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	DB = db
	return db
}

func InitializeDB() *pgxpool.Pool {
	db := connectToDB()
	return db
}
