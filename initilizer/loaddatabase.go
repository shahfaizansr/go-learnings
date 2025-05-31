package initilizer

import (
	"log"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func LoadDataBase() {
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")

	}
}
