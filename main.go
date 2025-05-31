package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shahfaizansr/initilizer"
	"github.com/shahfaizansr/service"
)

func init() {
	initilizer.LoadEnvFile()
	initilizer.LoadDataBase()
}

func main() {

	router := gin.Default()
	router.POST("/add", service.PostService)
	router.POST("/getall", service.GetAllService)
	router.GET("/post/:id", service.GetService)
	router.Run()

}
