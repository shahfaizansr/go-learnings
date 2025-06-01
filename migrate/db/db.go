package sqldb

import (
	"fmt"
	"log"
	"os"

	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// DBHandler encapsulates GORM DB instance
type DBHandler struct {
	DB     *gorm.DB
	Logger *logharbour.Logger
	Rigel  *rigel.Rigel
}

// NewSQLServerDB connects to a SQL Server DB using the given connection string
func NewSQLServerDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL") // Expected format: "sqlserver://username:password@localhost:1433?database=yourdb"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("SQL Server ping failed: %w", err)
	}

	log.Println("✅ Successfully connected to SQL Server database")
	return db, nil
}

// NewSQLServerHandler returns a DBHandler with optional logger and rigel client
func NewSQLServerHandler(logger *logharbour.Logger, rigelClient *rigel.Rigel) (*DBHandler, error) {
	db, err := NewSQLServerDB()
	if err != nil {
		return nil, err
	}

	return &DBHandler{
		DB:     db,
		Logger: logger,
		Rigel:  rigelClient,
	}, nil
}
