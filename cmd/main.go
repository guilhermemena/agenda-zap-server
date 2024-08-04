package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/guilhermemena/agenda-zap-server/cmd/api"
	"github.com/guilhermemena/agenda-zap-server/cmd/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), configs.Envs.DBConnection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), dbpool)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
