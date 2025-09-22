package main

import (
	"database/sql"
	"log"

	"github.com/Viczdera/bank-gorm/api"
	db "github.com/Viczdera/bank-gorm/db/sqlc"
	"github.com/Viczdera/bank-gorm/util"

	_ "github.com/lib/pq"
)

var (
	SERVER_NO_START = "server no gree start! 😭"
	CONFIG_NO_LOAD  = "Failed to load config 💿"
	DB_NO_CONNECT   = "Could not connect to DB 💿"
)

func main() {
	//load config from env
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(CONFIG_NO_LOAD, err)
	}

	//establish connection to db
	connDB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(DB_NO_CONNECT, err)
	}

	store := db.NewStore(connDB)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal(SERVER_NO_START, err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal(SERVER_NO_START, err)
	}

}
