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
		imgs.GET("/:image", s.BindUUID, s.handleImagesGet)
		imgs.DELETE("/:image", s.BindUUID, s.handleImagesDelete)
	}
}
