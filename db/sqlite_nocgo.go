//go:build !(cgo && (linux || darwin))

package db

import (
	"github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/nickfthedev/fiberHTMX/lib"
	"gorm.io/gorm"
)

func ConnectSQLite(config lib.ConfigStruct) (*gorm.DB, error) {
	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open(config.DbFileName), &gorm.Config{QueryFields: true})
	return db, err
}
