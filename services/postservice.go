package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shahfaizansr/initilizer"
	"github.com/shahfaizansr/models"
)

func PostService(ctx *gin.Context) {

	post := models.Post{
		Title: "Harry Potter",
		Body:  "Magical Movie",
	}

	result := initilizer.DB.Create(&post)

	if result.Error != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}
