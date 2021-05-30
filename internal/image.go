package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/minio/minio-go"
)

var _ ImageService = (*imageService)(nil)

var supportedContentTypes = []string{"image/jpeg", "image/png", "image/svg+xml"}

var (
	// ErrImageNotFound file not found
	ErrImageNotFound = errors.New("image not found")

	// ErrUnsupportedContentType file content type not supported. See supportedContentTypes
	ErrUnsupportedContentType = fmt.Errorf("unsupported content type: [%s]", strings.Join(supportedContentTypes, ", "))
)

type mustMakeDownloadURL func(string) string

// Image describes our base image type
type Image struct {
	Name        string    `json:"name"`
	Content     io.Reader `json:"-"`
	ContentType string    `json:"-"`
	Description string    `json:"description"`
	DownloadURL string    `json:"download_url"`
	Size        int64     `json:"-"`
}

func (i *Image) validateContentType() error {
	for _, c := range supportedContentTypes {
		if c == i.ContentType {
			return nil
		}
	}

	return ErrUnsupportedContentType
}

// newImage create a new Image
func newImage(name string, description string, header *multipart.FileHeader) (*Image, error) {
	// Read the file
	f, err := header.Open()
	if err != nil {
		return nil, err
	}

	// Default name is the filename
	var objectName = filepath.Base(header.Filename)

	// If user specified a custom name, use it and append the original file extension
	if name != "" {
		objectName = name + filepath.Ext(objectName)
	}

	return &Image{
		Name:        objectName,
		Description: description,
		DownloadURL: "",
		Content:     f,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
	}, nil
}

// ImageService describes available operations on Image
type ImageService interface {
	// Create make a new image
	Create(ctx context.Context, image *Image) (*Image, error)

	// Get return Image matching given uuid
	Get(ctx context.Context, name string) (*Image, error)

	// List returns a set of all images
	List(ctx context.Context) ([]*Image, error)

	// Delete deletes Image matching given uuid
	Delete(ctx context.Context, name string) error
}

type imageService struct {
	Minio      *minio.Client //  Minio is S3 compatible so we can safely use it
	BucketName string        // Bucket to work with
}

func (i *imageService) Create(ctx context.Context, image *Image) (*Image, error) {
	if err := image.validateContentType(); err != nil {
		return nil, err
	}

	_, err := i.Minio.PutObjectWithContext(ctx, i.BucketName, image.Name, image.Content, image.Size, minio.PutObjectOptions{
		ContentType: image.ContentType,
		UserMetadata: map[string]string{
			"description": image.Description,
		},
	})
	if err != nil {
		return nil, err
	}

	image.DownloadURL = i.mustMakeDownloadURL(image.Name)
	return image, nil
}

func (i *imageService) Get(ctx context.Context, name string) (*Image, error) {
	object, err := i.Minio.GetObject(i.BucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	info, err := object.Stat()
	if err != nil {
		e := minio.ToErrorResponse(err)
		switch e.StatusCode {
		case http.StatusNotFound:
			return nil, ErrImageNotFound
		default:
			return nil, err
		}
	}

	image := i.makeImage(&info)
	image.Content = object
	return image, nil
}

func (i *imageService) List(ctx context.Context) ([]*Image, error) {
	done := make(chan struct{})
	var images = make([]*Image, 0)

	go func() {
		defer close(done)
		for object := range i.Minio.ListObjectsV2(i.BucketName, "", false, done) {
			if err := object.Err; err != nil {
				log.Errorln(err)
				continue
			}

			// Get the real object for metadata
			image, err := i.Get(ctx, object.Key)
			if err == nil { // no error, use it
				images = append(images, image)
				continue
			}

			// error occurred, use the image without metadata
			images = append(images, i.makeImage(&object))
		}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return images, nil
	}
}

func (i *imageService) Delete(ctx context.Context, name string) error {
	return i.Minio.RemoveObject(i.BucketName, name)
}

func (i *imageService) makeImage(object *minio.ObjectInfo) *Image {
	return &Image{
		Name:        object.Key,
		Content:     nil,
		ContentType: object.ContentType,
		Description: object.Metadata.Get("X-Amz-Meta-Description"),
		DownloadURL: i.mustMakeDownloadURL(object.Key),
		Size:        object.Size,
	}
}

// mustMakeDownloadURL use the native Minio feature to generate download links
// each URLs will be available 7 days.
// TODO: better way ?
func (i *imageService) mustMakeDownloadURL(name string) string {
	d, _ := time.ParseDuration("604800s") // 7 days
	u, err := i.Minio.PresignedGetObject(i.BucketName, name, d, nil)
	if err != nil {
		log.Panicf("cannot generate presigned URL: %v", err)
	}

	return u.String()
}
