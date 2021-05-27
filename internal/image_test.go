package internal

import (
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"reflect"
	"testing"
	"time"

	"github.com/minio/minio-go"
)

func Test_imageService_makeImage(t *testing.T) {
	type fields struct {
		Minio                     *minio.Client
		BucketName                string
		GenerateDownloadRouteFunc GenerateDownloadRouteFunc
	}
	type args struct {
		object *minio.ObjectInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Image
	}{
		{
			name: "Expected",
			fields: fields{
				GenerateDownloadRouteFunc: func(s string) string {
					return s
				},
			},
			args: args{
				object: &minio.ObjectInfo{
					Key:         "foo",
					Size:        1234,
					ContentType: "my-content-type",
					Metadata: func() http.Header {
						h := make(http.Header)
						h.Set("X-Amz-Meta-Description", "foo")
						return h
					}(),
				},
			},
			want: &Image{
				Name:        "foo",
				Content:     nil,
				ContentType: "my-content-type",
				Description: "foo",
				DownloadURL: "foo",
				Size:        1234,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &imageService{
				Minio:                     tt.fields.Minio,
				BucketName:                tt.fields.BucketName,
				GenerateDownloadRouteFunc: tt.fields.GenerateDownloadRouteFunc,
			}
			if got := i.makeImage(tt.args.object); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newImageService(t *testing.T) {
	type args struct {
		minio                     *minio.Client
		bucketName                string
		generateDownloadRouteFunc GenerateDownloadRouteFunc
	}

	expectedClient := new(minio.Client)

	tests := []struct {
		name string
		args args
		want *imageService
	}{
		{
			name: "Expected",
			args: args{
				minio:                     expectedClient,
				bucketName:                "foo",
				generateDownloadRouteFunc: nil,
			},
			want: &imageService{
				Minio:                     expectedClient,
				BucketName:                "foo",
				GenerateDownloadRouteFunc: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newImageService(tt.args.minio, tt.args.bucketName, tt.args.generateDownloadRouteFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newImageService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newImage(t *testing.T) {
	type args struct {
		name        string
		description string
		file        *multipart.File
		header      *multipart.FileHeader
	}

	expectedFile := new(multipart.File)

	tests := []struct {
		name string
		args args
		want *Image
	}{
		{
			name: "Default name",
			args: args{
				name:        "",
				description: "foo",
				file:        expectedFile,
				header: &multipart.FileHeader{
					Filename: "foo.png",
					Header: func() textproto.MIMEHeader {
						h := make(textproto.MIMEHeader)
						h.Set("Content-Type", "foo/bar")
						return h
					}(),
					Size: 1234,
				},
			},
			want: &Image{
				Name:        "foo.png",
				Content:     nil,
				ContentType: "foo/bar",
				Description: "foo",
				DownloadURL: "",
				Size:        1234,
			},
		},
		{
			name: "Custom name with extension",
			args: args{
				name:        "hello.pdf",
				description: "foo",
				file:        expectedFile,
				header: &multipart.FileHeader{
					Filename: "foo.png",
					Header: func() textproto.MIMEHeader {
						h := make(textproto.MIMEHeader)
						h.Set("Content-Type", "foo/bar")
						return h
					}(),
					Size: 1234,
				},
			},
			want: &Image{
				Name:        "hello.pdf.png",
				Content:     nil,
				ContentType: "foo/bar",
				Description: "foo",
				DownloadURL: "",
				Size:        1234,
			},
		},
		{
			name: "Custom name without extension",
			args: args{
				name:        "hello",
				description: "foo",
				file:        expectedFile,
				header: &multipart.FileHeader{
					Filename: "foo.png",
					Header: func() textproto.MIMEHeader {
						h := make(textproto.MIMEHeader)
						h.Set("Content-Type", "foo/bar")
						return h
					}(),
					Size: 1234,
				},
			},
			want: &Image{
				Name:        "hello.png",
				Content:     nil,
				ContentType: "foo/bar",
				Description: "foo",
				DownloadURL: "",
				Size:        1234,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newImage(tt.args.name, tt.args.description, tt.args.file, tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImage_validateContentType(t *testing.T) {
	type fields struct {
		Name        string
		Content     io.Reader
		ContentType string
		Description string
		DownloadURL string
		Size        int64
	}

	rand.Seed(time.Now().Unix())
	var getRandomValidContentType = func() string {
		i := rand.Intn(len(supportedContentTypes) - 1)
		return supportedContentTypes[i]
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Valid",
			fields: fields{
				Name:        "",
				Content:     nil,
				ContentType: getRandomValidContentType(),
				Description: "",
				DownloadURL: "",
				Size:        0,
			},
			wantErr: false,
		},
		{
			name: "Invalid",
			fields: fields{
				Name:        "",
				Content:     nil,
				ContentType: "foo",
				Description: "",
				DownloadURL: "",
				Size:        0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Image{
				Name:        tt.fields.Name,
				Content:     tt.fields.Content,
				ContentType: tt.fields.ContentType,
				Description: tt.fields.Description,
				DownloadURL: tt.fields.DownloadURL,
				Size:        tt.fields.Size,
			}
			if err := i.validateContentType(); (err != nil) != tt.wantErr {
				t.Errorf("validateContentType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
