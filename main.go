package main

import (
	"database/sql"
	"log"
	"lost_found/api"
	"lost_found/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	store := api.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("failed to create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can not start server:", err)
	}
}
