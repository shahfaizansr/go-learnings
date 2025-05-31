package service

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shahfaizansr/initilizer"
	"github.com/shahfaizansr/models"
)

func GetService(ctx *gin.Context) {

	idParam := ctx.Param("id") // read :id from URL
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var post models.Post
	result := initilizer.DB.First(&post, id)

	if result.Error != nil {

		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	ctx.JSON(http.StatusOK, post)
}
