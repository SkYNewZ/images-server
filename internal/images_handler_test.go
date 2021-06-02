package internal

import (
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

func Test_server_handleImagesDelete(t *testing.T) {
	type fields struct {
		image uuid.UUID
		Image ImageService
	}

	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "Image does not exist",
			fields: fields{
				image: uuid.New(),
				Image: new(testingImageService), // fresh service, this image does not exist
			},
			want:    500,
			wantErr: true,
		},
		{
			name: "Image exist",
			fields: func() fields {
				f := new(fields)
				f.image = uuid.New()
				f.Image = new(testingImageService)

				// Make an image
				_, _ = f.Image.Create(context.TODO(), &Image{
					Key:         f.image,
					Name:        "",
					Content:     nil,
					ContentType: "image/png",
					Description: "",
					DownloadURL: "",
					Size:        0,
				})

				return *f
			}(),
			want:    204,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				router: gin.New(),
				Image:  tt.fields.Image,
			}

			var err error
			// basic middleware to read received error
			readErrorHandler := func(c *gin.Context) {
				c.Next()
				errs := c.Errors.ByType(gin.ErrorTypeAny)

				if len(errs) == 0 {
					return
				}

				err = errs[0]
				c.Status(500)
			}

			// Make the fake route
			s.router.GET("/foo/:image", readErrorHandler, s.BindUUID, s.handleImagesDelete)

			// Run test
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/foo/"+tt.fields.image.String(), nil)
			s.ServeHTTP(rw, req)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, rw.Code)
		})
	}
}

func Test_server_handleImagesList(t *testing.T) {
	type fields struct {
		Image ImageService
	}

	type args struct {
		ctx context.Context
	}

	s := new(testingImageService)

	// Fill images buffer
	count := 10
	var images = make([]*Image, count)
	for i := range images {
		images[i] = &Image{
			Key:         uuid.New(),
			Name:        "foo",
			Content:     nil,
			ContentType: "image/png",
			Description: "bar",
			DownloadURL: "",
			Size:        1234,
		}

		_, _ = s.Create(context.TODO(), images[i])
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		want     string
		wantCode int
	}{
		{
			name:   "Expected",
			fields: fields{s},
			args: args{
				ctx: context.Background(),
			},
			want:     "",
			wantCode: 200,
		},
		{
			name:   "Error",
			fields: fields{s},
			args: args{
				ctx: context.WithValue(context.Background(), "error", true), //nolint:revive
			},
			want:     "",
			wantCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				router: gin.New(),
				Image:  tt.fields.Image,
			}

			// basic middleware to read received error
			readErrorHandler := func(c *gin.Context) {
				c.Next()
				errs := c.Errors.ByType(gin.ErrorTypeAny)

				if len(errs) == 0 {
					return
				}

				c.Status(500)
			}

			s.router.GET("/foo", readErrorHandler, s.handleImagesList)
			rw := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/foo", nil)
			req = req.WithContext(tt.args.ctx)
			s.ServeHTTP(rw, req)

			assert.Equal(t, tt.wantCode, rw.Code)
			//assert.JSONEq(t, tt.want, rw.Body.String())
		})
	}
}
