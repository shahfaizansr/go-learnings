package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shahfaizansr/initilizer"
	"github.com/shahfaizansr/models"
)

func GetAllService(ctx *gin.Context) {
	var posts []models.Post

	result := initilizer.DB.Find(&posts)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}
