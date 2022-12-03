package db

import (
	"database/sql"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/escalopa/gobank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMain(m *testing.M) {
	var err error
	// Load config
	config := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load configuration for testing", err)
	}

	// Connect to db
	testDB, err = sql.Open(config.Driver, config.ConnectionString)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// Set connection & run tests
	testQueries = New(testDB)
	os.Exit(m.Run())
}
