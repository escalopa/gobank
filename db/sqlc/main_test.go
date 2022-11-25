package db

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/escalopa/go-bank/utils"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMain(m *testing.M) {
	var err error
	var config utils.Config
	// Load config
	config, err = utils.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load configuration for testing", err)
	}

	// Connect to db
	testDB, err = sql.Open(config.DB.Driver, config.DB.ConnectionString)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// Set connection & run tests
	testQueries = New(testDB)
	os.Exit(m.Run())
}
