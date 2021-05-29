package internal

func (s *server) routes() {
	// TODO: Use API versioning ?

	// Health check
	s.router.GET("/_health", s.handleHealthCheck)

	// Images
	imgs := s.router.Group("/images")
	imgs.Use(s.handleErrors)
	{
		imgs.GET("", s.handleImagesList)
		imgs.POST("", s.handleImagesCreate)
		imgs.GET("/:image", s.handleImagesGet)
		imgs.DELETE("/:image", s.handleImagesDelete)
	}

	// Download
	s.router.GET("/download/:image", s.handleErrors, s.handleImagesDownload)
}
