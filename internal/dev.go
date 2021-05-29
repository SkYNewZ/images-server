// +build dev
//go:build dev

package internal

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Use 'go run -tags dev .' to run the server in debug mode
func init() {
	ginMode = gin.DebugMode
	log.SetLevel(log.TraceLevel)
}
