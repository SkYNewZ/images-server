package internal

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type ContextKey int

const imageIDContextKey ContextKey = iota

// handleErrors carry our HTTP handlers error logic
func (s *server) handleErrors(handler HandlerFuncWithError) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := handler(w, req)
		if err == nil {
			return
		}

		log.Errorln(err)

		// Default returned error
		var e = &Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}

		// Is it our custom error ?
		if errors.As(err, &e) {
			e = err.(*Error)
		}

		_ = s.respond(w, req, e, e.Code)
	}
}

// respond send a JSON response
func (s *server) respond(w http.ResponseWriter, _ *http.Request, data interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	if data != nil {
		return json.NewEncoder(w).Encode(data)
	}

	return nil
}

// imageIDFromContext return image name stored in given context
func imageIDFromContext(ctx context.Context) string {
	return ctx.Value(imageIDContextKey).(string)
}

// handleGetImageName read current http.Request params and set image's name in the request context
func (s *server) handleGetImageName(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		*r = *r.WithContext(context.WithValue(r.Context(), imageIDContextKey, vars["image"]))
		next.ServeHTTP(w, r)
	}
}
