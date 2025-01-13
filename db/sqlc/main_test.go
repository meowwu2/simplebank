package db

import (
	"database/sql"
	"log"
	"os"
	"simplebank/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries


var testDB *sql.DB
func TestMain(m *testing.M) {
	config,err:=util.LoadConfig("../..")
	if err!=nil{
		log.Fatal("load config err: ",err)
	}
	testDB,err =sql.Open(config.DBDriver,config.DBSource)
	if err!=nil{
		log.Fatal("cannot connect database:",err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}