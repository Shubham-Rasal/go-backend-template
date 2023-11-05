package main

import (
	"database/sql"
	"log"

	"github.com/Shubham-Rasal/blog-backend/api"
	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:password@localhost:5432/blog?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	DB, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(DB)
	server := api.NewServer(*store)
	err = server.Start(serverAddress)

}
