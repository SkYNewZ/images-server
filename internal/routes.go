package internal

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *server) routes() {
	//var uuidPattern = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

	// GET /_health
	s.router.HandleFunc("/_health", s.handleErrors(s.handleHealthCheck)).
		Methods(http.MethodGet)

	// GET /images - list images
	s.router.HandleFunc("/images", s.handleErrors(s.handleImagesList)).
		Methods(http.MethodGet)

	// POST /images - create an image
	s.router.HandleFunc("/images", s.handleErrors(s.handleImagesCreate)).
		Methods(http.MethodPost)

	// GET /images/<image name> - get one image
	s.router.HandleFunc("/images/{image}", s.handleGetImageName(s.handleErrors(s.handleImagesGet))).
		Methods(http.MethodGet)

	// DELETE /images/<image name> - delete an image
	s.router.HandleFunc("/images/{image}", s.handleGetImageName(s.handleErrors(s.handleImagesDelete))).
		Methods(http.MethodDelete)

	// GET /download/<image name> - download given image
	s.router.HandleFunc("/download/{image}", s.handleGetImageName(s.handleErrors(s.handleImagesDownload))).
		Methods(http.MethodGet).Name("download")
}

func (s *server) GenerateDownloadURL(domain string) GenerateDownloadRouteFunc {
	return func(name string) string {
		u, err := s.router.Get("download").URL("image", name)
		if err != nil {
			log.Errorln(err)
			return ""
		}

		return domain + u.String()
	}
}
