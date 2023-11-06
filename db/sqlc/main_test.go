package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Shubham-Rasal/blog-backend/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	//load config
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load env", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())

}
