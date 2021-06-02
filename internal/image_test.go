package internal

import (
	"io"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

// makeFileHeader generate a valid multipart.FileHeader with the given file
// https://stackoverflow.com/a/57084200
func makeFileHeader(t *testing.T, name string) (multipart.File, *multipart.FileHeader) {
	t.Helper()

	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	throw := func(err error) {
		t.Fatal(err)
	}

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("file", name)
		if err != nil {
			throw(err)
			return
		}

		f, err := os.Open("testdata/" + name)
		if err != nil {
			throw(err)
			return
		}

		if _, err = io.Copy(part, f); err != nil {
			throw(err)
			return
		}
	}()

	request := httptest.NewRequest("POST", "/foo", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	file, header, err := request.FormFile("file")
	if err != nil {
		throw(err)
	}

	return file, header
}

// TestImage_validateContentType is more like a non-regression test
func TestImage_validateContentType(t *testing.T) {
	type fields struct {
		Key         uuid.UUID
		Name        string
		Content     io.Reader
		ContentType string
		Description string
		DownloadURL string
		Size        int64
	}

	type T struct {
		name    string
		fields  fields
		wantErr bool
	}

	invalidContentTypes := []string{
		"image/x-icon",
		"image/bmp",
		"image/tiff",
		"application/zip",
		"video/mpeg",
	}

	validContentTypes := []string{
		"image/jpeg",
		"image/png",
		"image/svg+xml",
	}

	// Populate tests
	var tests = make([]T, 0)
	for _, s := range invalidContentTypes {
		tests = append(tests, T{
			name: "Invalid: " + s,
			fields: fields{
				Key:         uuid.UUID{},
				Name:        "",
				Content:     nil,
				ContentType: s,
				Description: "",
				DownloadURL: "",
				Size:        0,
			},
			wantErr: true,
		})
	}

	for _, s := range validContentTypes {
		tests = append(tests, T{
			name: "Valid: " + s,
			fields: fields{
				Key:         uuid.UUID{},
				Name:        "",
				Content:     nil,
				ContentType: s,
				Description: "",
				DownloadURL: "",
				Size:        0,
			},
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Image{
				Key:         tt.fields.Key,
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

func Test_newImage(t *testing.T) {
	type args struct {
		name        string
		description string
		header      *multipart.FileHeader
	}

	file, header := makeFileHeader(t, "gopher.png")

	tests := []struct {
		name    string
		args    args
		want    *Image
		wantErr bool
	}{
		{
			name: "Default name",
			args: args{
				name:        "",
				description: "",
				header:      header,
			},
			want: &Image{
				Key:         uuid.UUID{},
				Name:        "gopher.png",
				Content:     file,
				ContentType: "application/octet-stream",
				Description: "",
				DownloadURL: "",
				Size:        254145,
			},
			wantErr: false,
		},
		{
			name: "Default name with description",
			args: args{
				name:        "",
				description: "foo",
				header:      header,
			},
			want: &Image{
				Key:         uuid.UUID{},
				Name:        "gopher.png",
				Content:     file,
				ContentType: "application/octet-stream",
				Description: "foo",
				DownloadURL: "",
				Size:        254145,
			},
			wantErr: false,
		},
		{
			name: "Custom name",
			args: args{
				name:        "hello",
				description: "foo",
				header:      header,
			},
			want: &Image{
				Key:         uuid.UUID{},
				Name:        "hello.png",
				Content:     file,
				ContentType: "application/octet-stream",
				Description: "foo",
				DownloadURL: "",
				Size:        254145,
			},
			wantErr: false,
		},
		{
			name: "Custom name but original extension",
			args: args{
				name:        "hello.pdf",
				description: "foo",
				header:      header,
			},
			want: &Image{
				Key:         uuid.UUID{},
				Name:        "hello.pdf.png",
				Content:     file,
				ContentType: "application/octet-stream",
				Description: "foo",
				DownloadURL: "",
				Size:        254145,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newImage(tt.args.name, tt.args.description, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("newImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Key is randomly generated
			// Don't compare the reader
			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreFields(Image{}, "Key"), cmpopts.IgnoreUnexported(io.SectionReader{})); diff != "" {
				t.Errorf("newImage() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
