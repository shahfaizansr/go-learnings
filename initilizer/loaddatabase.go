package initilizer

import (
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

func LoadDataBase() {
	dsn := os.Getenv("DB_URL")
	// DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("failed to connect database:", err)
	// }
	// log.Println("✅ Migration completed successfully!")

	// // Auto migrate table
	// if err := DB.AutoMigrate(&mycalc.CalcResponseModel{}); err != nil {
	// 	log.Fatal("failed to auto migrate:", err)
	// }

	db, err := sqlx.Connect("sqlserver", dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	log.Println("✅ Successfully connected to SQL Server using sqlx")
}
