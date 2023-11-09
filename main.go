package main

import (
	"database/sql"
	"log"

	"github.com/Shubham-Rasal/blog-backend/api"
	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/Shubham-Rasal/blog-backend/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
)

func main() {

	//load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env", err)
	}
	log.Println("loaded config")

	DB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(DB)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server", err)
	}
	err = server.Start(config.ServerAddress)
}
