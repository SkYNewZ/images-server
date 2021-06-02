package internal

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

var (
	_      ImageService = (*testingImageService)(nil)
	images              = make(map[uuid.UUID]*Image)
)

type testingImageService struct{}

func (t *testingImageService) Create(ctx context.Context, image *Image) (*Image, error) {
	if image.ContentType != "image/png" {
		// just to throw an error on create
		return nil, ErrUnsupportedContentType
	}

	// Fake download URL
	image.DownloadURL = "https://example.com/download/" + image.Key.String()
	images[image.Key] = image

	return image, nil
}

func (t *testingImageService) Get(ctx context.Context, id uuid.UUID) (*Image, error) {
	if v, ok := images[id]; ok {
		return v, nil
	}

	return nil, ErrImageNotFound
}

func (t *testingImageService) List(ctx context.Context) ([]*Image, error) {
	if _, ok := ctx.Value("error").(bool); ok {
		return nil, fmt.Errorf("oops")
	}

	res := make([]*Image, 0)
	for _, v := range images {
		res = append(res, v)
	}

	return res, nil
}

func (t *testingImageService) Delete(ctx context.Context, ids ...uuid.UUID) error {
	for _, id := range ids {
		if _, ok := images[id]; !ok {
			return fmt.Errorf("oops")
		}

		delete(images, id)
	}

	return nil
}
