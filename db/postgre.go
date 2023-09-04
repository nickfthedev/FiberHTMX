package db

import (
	"fmt"
	"log"

	"github.com/nickfthedev/fiberHTMX/lib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgre(config lib.ConfigStruct) (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	var dsn string
	if config.PostgresURL != "" {
		dsn = config.PostgresURL
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			config.DbHost,
			config.DbUser,
			config.DbPass,
			config.DbName,
			config.DbPort,
		)
	}
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{QueryFields: true})
	if err != nil {
		log.Println(err)
		panic("Failled to connect to Database. ")
	}
	return db, nil

}
