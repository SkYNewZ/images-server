package internal

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

var (
	// ErrCannotFindFile describes error when request form does not contains the 'file' key
	ErrCannotFindFile = errors.New("cannot find 'file'")
)

// uploadImageForm describes expected request form to properly upload an image
// https://pkg.go.dev/github.com/go-playground/validator?utm_source=godoc#hdr-Baked_In_Validators_and_Tags
type uploadImageForm struct {
	Name        string                `form:"name" binding:"-"`
	Description string                `form:"description" binding:"-"`
	Header      *multipart.FileHeader `form:"file" binding:"required"`
}

func (s *server) handleImagesList(c *gin.Context) {
	images, err := s.Image.List(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, images)
}

func (s *server) handleImagesGet(c *gin.Context) {

	id, _ := c.Get(UUIDContextKey)
	image, err := s.Image.Get(c.Request.Context(), id.(uuid.UUID))
	if err != nil {
		if errors.Is(err, ErrImageNotFound) {
			err = newNotFoundError(err)
		}

		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, image)
}

func (s *server) handleImagesCreate(c *gin.Context) {
	var form uploadImageForm
	if err := c.ShouldBindWith(&form, binding.FormMultipart); err != nil {
		_ = c.Error(newBadRequestError(err))
		return
	}

	// create image object
	image, err := newImage(form.Name, form.Description, form.Header)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			err = newBadRequestError(ErrCannotFindFile)
		}

		_ = c.Error(err)
		return
	}

	// upload it!
	image, err = s.Image.Create(c.Request.Context(), image)
	if err != nil {
		if errors.Is(err, ErrUnsupportedContentType) {
			err = newUnsupportedMediaType(err)
		}

		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, image)
}

func (s *server) handleImagesDelete(c *gin.Context) {
	id, _ := c.Get(UUIDContextKey)
	if err := s.Image.Delete(c.Request.Context(), id.(uuid.UUID)); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
