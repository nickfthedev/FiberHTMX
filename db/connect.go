package db

import (
	"fmt"

	"github.com/nickfthedev/fiberHTMX/lib"
	"gorm.io/gorm"
)

// Database instance
var DB *gorm.DB

// Connect to database and put instance into var DB
func ConnectDB(config lib.ConfigStruct) {
	var err error

	switch config.DbDriver {
	case "POSTGRESQL":
		fmt.Println("DB Driver: PostgreSQL")
		DB, err = ConnectPostgre(config)
		if err != nil {
			panic(err)
		}

	case "SQLITE":
		fmt.Println("DB Driver: SQLite")
		DB, err = ConnectSQLite(config)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println("Check your Databasedriver!")
		panic("no db driver")
	}

}
