package main

import (
	"database/sql"
	"log"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/services"
	"github.com/dato7898/grpc-tube/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := services.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start()
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
