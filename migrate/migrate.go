package main

import (
	"github.com/shahfaizansr/initilizer"
)

func init() {
	initilizer.LoadDataBase()
	initilizer.LoadEnvFile()
}

// func main() {

// 	// Migrate the schema
// 	initilizer.DB.AutoMigrate(&models.Post{})
// }
