package main

import (
	"database/sql"
	"log"

	db "github.com/dato7898/grpc-tube/db/sqlc"
	"github.com/dato7898/grpc-tube/services"

	_ "github.com/lib/pq"
)

func main() {

	conn, err := sql.Open("postgres", "postgresql://grpc_tube:grpc_tube@grpc-tube-db:5432/grpc_tube?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := services.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start()
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
