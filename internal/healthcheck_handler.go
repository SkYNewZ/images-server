package internal

import (
	"net/http"
)

func (s *server) handleHealthCheck(w http.ResponseWriter, r *http.Request) error {
	data := map[string]interface{}{
		"ok":      true,
		"version": buildNumber,
	}
	return s.respond(w, r, data, http.StatusOK)
}
