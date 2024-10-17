package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB(DBPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error opening sqlite database: %v", err)
		return nil
	}
	return db
}
