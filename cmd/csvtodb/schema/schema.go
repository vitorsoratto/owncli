package schema

import (
	"owncli/cmd/csvtodb/database"

	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(DBPath string) {
	db = database.GetDB(DBPath)
}
