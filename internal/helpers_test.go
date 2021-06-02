package internal

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"
)

func Test_server_handleErrors(t *testing.T) {
	type args struct {
		handler gin.HandlerFunc
	}
	tests := []struct {
		name     string
		args     args
		want     string
		wantCode int
	}{
		{
			name: "No error",
			args: args{
				handler: func(c *gin.Context) {
					c.String(200, "foo")
				},
			},
			want:     "foo",
			wantCode: 200,
		},
		{
			name: "Generic error",
			args: args{
				handler: func(c *gin.Context) {
					_ = c.Error(fmt.Errorf("oops"))
				},
			},
			want:     `{"code":500,"message":"oops"}`,
			wantCode: 500,
		},
		{
			name: "Custom error",
			args: args{
				handler: func(c *gin.Context) {
					err := new(Error)
					err.Code = 400
					err.Message = "argh"
					_ = c.Error(err)
				},
			},
			want:     `{"code":400,"message":"argh"}`,
			wantCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := new(server)
			s.router = gin.New()
			s.router.GET("/foo", s.handleErrors, tt.args.handler)

			rw := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/foo", nil)
			s.ServeHTTP(rw, req)

			assert.Equal(t, rw.Code, tt.wantCode)
			assert.Equal(t, rw.Body.String(), tt.want)
		})
	}
}
