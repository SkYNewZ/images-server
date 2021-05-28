package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

var _ http.Handler = (*server)(nil)

type HandlerFuncWithError func(http.ResponseWriter, *http.Request) error

const (
	bucketNameImages = "images"
)

// server handle our global server instance logic
// Inspired by https://youtu.be/rWBSMsLG8po?t=613
type server struct {
	router *mux.Router
	Image  ImageService
}

// ServeHTTP implements http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}

func newServer(domain string) *server {
	s := new(server)
	s.router = mux.NewRouter()
	s.routes()      // declare our routes
	s.middlewares() // declare our middlewares

	// Inject dependencies
	s.Image = newImageService(newMinio(bucketNameImages), bucketNameImages, s.GenerateDownloadURL(domain))

	return s
}
