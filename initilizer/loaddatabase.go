package initilizer

import (
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

func LoadDataBase() {
	dsn := os.Getenv("DB_URL")
	db, err := sqlx.Connect("sqlserver", dsn)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	log.Println("✅ Successfully connected to SQL Server using sqlx")
}
