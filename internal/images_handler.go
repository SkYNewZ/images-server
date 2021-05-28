package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	// ErrCannotFindFile describes error when request form does not contains the 'file' key
	ErrCannotFindFile = errors.New("cannot find 'file'")
)

func (s *server) handleImagesList(w http.ResponseWriter, req *http.Request) error {
	images, err := s.Image.List(req.Context())
	if err != nil {
		return err
	}

	return s.respond(w, req, images, http.StatusOK)
}

func (s *server) handleImagesGet(w http.ResponseWriter, req *http.Request) error {
	imageUUID := imageIDFromContext(req.Context())
	image, err := s.Image.Get(req.Context(), imageUUID)
	if err != nil {
		switch {
		case errors.Is(err, ErrImageNotFound):
			return newNotFoundError(err)
		default:
			return err
		}
	}

	return s.respond(w, req, image, http.StatusOK)
}

func (s *server) handleImagesCreate(w http.ResponseWriter, req *http.Request) error {
	if err := req.ParseMultipartForm(10 << 20); err != nil { // 10MB
		switch {
		case errors.Is(err, http.ErrNotMultipart):
			return newBadRequestError(err)
		default:
			return err
		}
	}

	f, header, err := req.FormFile("file")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrMissingFile):
			return newBadRequestError(ErrCannotFindFile)
		default:
			return err
		}
	}
	defer f.Close()

	image, err := s.Image.Create(req.Context(), newImage(req.PostFormValue("name"), req.PostFormValue("description"), &f, header))
	if err != nil {
		switch {
		case errors.Is(err, ErrUnsupportedContentType):
			return newUnsupportedMediaType(err)
		default:
			return err
		}
	}

	return s.respond(w, req, image, http.StatusCreated)
}

func (s *server) handleImagesDelete(w http.ResponseWriter, req *http.Request) error {
	imageUUID := imageIDFromContext(req.Context())
	if err := s.Image.Delete(req.Context(), imageUUID); err != nil {
		return err
	}

	return s.respond(w, req, nil, http.StatusNoContent)
}

func (s *server) handleImagesDownload(w http.ResponseWriter, req *http.Request) error {
	imageUUID := imageIDFromContext(req.Context())
	image, err := s.Image.Get(req.Context(), imageUUID)
	if err != nil {
		switch {
		case errors.Is(err, ErrImageNotFound):
			return newNotFoundError(err)
		default:
			return err
		}
	}

	// https://stackoverflow.com/a/24116517
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", image.Name))
	w.Header().Set("Content-Type", image.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(image.Size, 10))
	_, err = io.Copy(w, image.Content)
	return err
}
