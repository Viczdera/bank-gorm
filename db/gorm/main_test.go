package gorm

import (
	"log"
	"os"
	"testing"

	"github.com/Viczdera/bank-gorm/util"
)

var testStore Store

// load config -> connect via DSN -> run migrations -> create store instance -> run tests

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../../")
	if err != nil {
		log.Fatal("\n 🔴 Failed to load config 💿: ", err)
	}

	//connecting via the DSN from config
	testDB, err := ConnectFromDSN(config.DBSource)
	if err != nil {
		log.Fatal("\n 🔴 Could not connect to DB: ", err)
	}

	if err := AutoMigrate(testDB); err != nil {
		log.Fatal("\n 🔴 Failed to run test migrations: ", err)
	}

	testStore = NewStore(testDB)

	os.Exit(m.Run())
}
