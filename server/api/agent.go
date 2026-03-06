package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListAgent(c *gin.Context) {

	c.Status(http.StatusOK)
}
