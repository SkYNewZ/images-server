package internal

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// loggingMiddleware enables http access log
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":      r.Method,
			"request_uri": r.RequestURI,
			"user_agent":  r.UserAgent(),
		}).Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// responseContentTypeMiddleware defines JSON as default returned Content-Type
func responseContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func (s *server) middlewares() {
	s.router.Use(loggingMiddleware)
	s.router.Use(responseContentTypeMiddleware)
}
