package utils

import (
	"whale/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB1 *gorm.DB

var DB2 *gorm.DB

var DB *gorm.DB

func ConnectDatabase(cfg config.Config) {

	db1, err := gorm.Open(mysql.Open(cfg.DBSource1), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	DB1 = db1

	if cfg.DBSource2 != "" {
		db2, err := gorm.Open(mysql.Open(cfg.DBSource2), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to database!")
		}
		DB2 = db2
	}

	DB = db1
}
