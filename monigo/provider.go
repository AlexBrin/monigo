package monigo

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	DB *gorm.DB
)

func DbConnect() error {
	var err error
	DB, err = gorm.Open("sqlite3", Config.DB.File)
	if err != nil {
		LogError("SQL open err: %s", err)
		return err
	}

	migrate()

	return nil
}

func migrate() {
	DB.AutoMigrate(&History{})
}