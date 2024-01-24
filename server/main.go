package main

import (
	"database/sql"
	"github.com/bysergr/simple-bank/api"
	db "github.com/bysergr/simple-bank/db/sqlc"
	"github.com/bysergr/simple-bank/token"
	"github.com/bysergr/simple-bank/utils"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	connection, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	store := db.NewStoreSQl(connection)
	server := api.NewServer(store, maker, config)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
