package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_server_handleHealthCheck(t *testing.T) {
	s := new(server)
	s.router = gin.New()
	s.routes()
	buildNumber = "foo"

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/_health", nil)
	s.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
	assert.Equal(t, `{"ok":true,"version":"foo"}`, rr.Body.String())
}
