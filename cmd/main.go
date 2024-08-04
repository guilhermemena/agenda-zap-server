package main

import (
	"fmt"
	"log"

	"github.com/guilhermemena/agenda-zap-server/cmd/api"
	"github.com/guilhermemena/agenda-zap-server/cmd/configs"
	"github.com/guilhermemena/agenda-zap-server/storage"
)

func main() {
	db := storage.InitializeDB()

	server := api.NewAPIServer(fmt.Sprintf(":%s", configs.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
