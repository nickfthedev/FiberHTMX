package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/nickfthedev/fiberHTMX/lib"
)

func TestConnectDB(t *testing.T) {
	//Create Config Struct
	c := lib.ConfigStruct{
		DbDriver:    "SQLITE",
		DbFileName:  "../tmp/test.db",
		PostgresURL: "postgres://pwdrsigk:lGrKfo8Et9y2fA7eM26UT5pd0s9YLNAq@mel.db.elephantsql.com/pwdrsigk",
		DbHost:      "test",
		DbName:      "test",
		DbUser:      "test",
		DbPass:      "test",
		DbPort:      "8888",
	}

	os.Remove(c.DbFileName)
	fmt.Println("tmp/test.db has been removed")
	// Test SQLite in Folder tests/test.db
	ConnectDB(c)
	fmt.Println("tmp/test.db has been created & connected")
	ConnectDB(c)
	fmt.Println("tmp/test.db has connected")

	// Test PostgreSQL
	c.DbDriver = "POSTGRESQL"
	ConnectDB(c)
	fmt.Println("Connection test to postgresql db succeded")
	//
}
