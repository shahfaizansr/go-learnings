package main

import (
	"github.com/shahfaizansr/initilizer"
	"github.com/shahfaizansr/models"
)

func init() {
	initilizer.LoadDataBase()
	initilizer.LoadEnvFile()
}

func main() {
	// Migrate the schema
	initilizer.DB.AutoMigrate(&models.Post{})
}
