package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthStatusCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
