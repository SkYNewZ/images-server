package internal

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_server_BindUUID(t *testing.T) {
	tests := []struct {
		name     string
		req      *http.Request
		want     uuid.UUID
		wantCode int
	}{
		{
			name:     "Invalid UUID",
			req:      httptest.NewRequest("GET", "/foo/bar", nil),
			want:     uuid.UUID{},
			wantCode: 404,
		},
		{
			name:     "Valid",
			req:      httptest.NewRequest("GET", "/foo/62c1a546-62b6-4cbc-bffb-03919b17dc3a", nil),
			want:     uuid.MustParse("62c1a546-62b6-4cbc-bffb-03919b17dc3a"),
			wantCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := new(server)
			s.router = gin.New()

			// make a dummy route to test our middleware
			s.router.GET("/foo/:image", s.BindUUID, func(c *gin.Context) {
				v, ok := c.Get(UUIDContextKey)
				assert.Equal(t, ok, true)
				assert.Equal(t, v, tt.want)
			})

			rr := httptest.NewRecorder()
			s.ServeHTTP(rr, tt.req)

			assert.Equal(t, rr.Code, tt.wantCode)
		})
	}
}

func Test_server_Recover(t *testing.T) {
	log.SetOutput(io.Discard)

	s := new(server)
	s.router = gin.New()
	s.router.Use(gin.CustomRecoveryWithWriter(io.Discard, s.Recover))

	// Route to trigger panic with custom string
	s.router.GET("/foo", func(c *gin.Context) {
		panic("oops")
	})

	// Route to trigger panic with default string
	s.router.GET("/bar", func(c *gin.Context) {
		panic(map[string]string{
			"foo": "bar",
		})
	})

	tests := []struct {
		name     string
		req      *http.Request
		want     string
		wantCode int
	}{
		{
			name:     "Custom message",
			req:      httptest.NewRequest("GET", "/foo", nil),
			want:     `{"code":500,"message":"oops"}`,
			wantCode: 500,
		},
		{
			name:     "Default message",
			req:      httptest.NewRequest("GET", "/bar", nil),
			want:     `{"code":500,"message":"Unexpected error occurred. Check logs for more details"}`,
			wantCode: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			s.ServeHTTP(rr, tt.req)

			assert.Equal(t, rr.Code, tt.wantCode)
			assert.Equal(t, rr.Body.String(), tt.want)
		})
	}
}
