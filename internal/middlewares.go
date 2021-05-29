package internal

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func (s *server) Recover(c *gin.Context, err interface{}) {
	e := new(Error)
	e.Code = http.StatusInternalServerError
	e.Message = "Unexpected error occurred. Check logs for more details"

	if err, ok := err.(string); ok {
		e.Message = err
	}

	log.Errorln(err)
	c.AbortWithStatusJSON(e.Code, e)
}

func (s *server) middlewares() {
	s.router.Use(gin.Logger())
	s.router.Use(gin.CustomRecovery(s.Recover))
	// TODO: use custom logger ?
}
