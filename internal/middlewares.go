package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

const UUIDContextKey = "uuid"

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

// BindUUID tries to read uuid.UUID from route param and set the value into context
// If param is not a valid uuid.UUID, an HTTP 404 response will be send
func (s *server) BindUUID(c *gin.Context) {
	type r struct {
		ID string `uri:"image" binding:"uuid4_rfc4122,required"`
	}

	var id r
	if err := c.ShouldBindUri(&id); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Set(UUIDContextKey, uuid.MustParse(id.ID))
}

func (s *server) middlewares() {
	s.router.Use(gin.Logger())
	s.router.Use(gin.CustomRecovery(s.Recover))
	// TODO: use custom logger ?
}
