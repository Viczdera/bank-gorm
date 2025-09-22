package db

//same package as crud code
//main_test used to set up connection with query object

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/Viczdera/bank-gorm/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

// testing
func TestMain(t *testing.M) {

	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("Failed to load config 💿", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	//use connection to create new test queries object
	testQueries = New(testDB)

	//run and report back to test runner via the o.exit command
	os.Exit(t.Run())

}
