package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *server) handleErrors(c *gin.Context) {
	c.Next()
	detectedErrors := c.Errors.ByType(gin.ErrorTypeAny)

	if len(detectedErrors) == 0 {
		return
	}

	err := detectedErrors[0].Err
	var e *Error
	switch v := err.(type) {
	case *Error:
		e = v
	default:
		e = &Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	c.JSON(e.Code, e)
}
