package initilizer

import (
	"log"
	"os"

	"github.com/shahfaizansr/models/mycalc"
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
		log.Fatal("failed to connect database:", err)
	}
	log.Println("✅ Migration completed successfully!")

	// Auto migrate table
	if err := DB.AutoMigrate(&mycalc.CalcResponseModel{}); err != nil {
		log.Fatal("failed to auto migrate:", err)
	}
}
