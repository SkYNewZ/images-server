package internal

import (
	"net/http"

	"github.com/SkYNewZ/images-server/internal/minio"
	"github.com/gin-gonic/gin"
)

var _ http.Handler = (*server)(nil)

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) error

const (
	bucketNameImages = "images"
)

var ginMode = gin.ReleaseMode

// server handle our global server instance logic
// Inspired by https://youtu.be/rWBSMsLG8po?t=613
type server struct {
	router *gin.Engine
	Image  ImageService
}

// ServeHTTP implements http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func newServer(domain string) *server {
	// Set default gin mode to release
	gin.SetMode(ginMode)

	s := new(server)
	s.router = gin.New()
	s.router.HandleMethodNotAllowed = true
	s.router.MaxMultipartMemory = 5 << 20 // 5MB

	s.routes()      // declare our routes
	s.middlewares() // declare our middlewares

	// Inject dependencies
	s.Image = &imageService{
		Minio: minio.New(bucketNameImages),
	}

	return s
}
