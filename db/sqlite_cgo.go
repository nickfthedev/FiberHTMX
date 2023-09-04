//go:build (darwin && cgo) || linux

package db

import (
	"github.com/nickfthedev/fiberHTMX/lib"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
)

func ConnectSQLite(config lib.ConfigStruct) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DbFileName), &gorm.Config{QueryFields: true})
	return db, err
}
