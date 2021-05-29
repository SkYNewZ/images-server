package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) handleHealthCheck(c *gin.Context) {
	data := map[string]interface{}{
		"ok":      true,
		"version": buildNumber,
	}

	c.JSON(http.StatusOK, data)
}
