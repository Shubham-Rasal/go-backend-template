package main

import (
	"database/sql"
	"log"

	"github.com/Shubham-Rasal/blog-backend/api"
	db "github.com/Shubham-Rasal/blog-backend/db/sqlc"
	"github.com/Shubham-Rasal/blog-backend/util"
	_ "github.com/lib/pq"
	_ "github.com/golang/mock/mockgen/model"
)

func main() {

	//load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load env", err)
	}

	DB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	store := db.NewStore(DB)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)

}
