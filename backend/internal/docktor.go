package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDocktor(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "I am a docktor!"})
}
