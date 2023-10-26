package pkg

import (
	"hack/internal/models"
	"hack/internal/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"

	pkgErrors "github.com/pkg/errors"
)

func ReadImage(file multipart.File, header *multipart.FileHeader) (*models.Image, error) {
	img := models.Image{
		Data: make([]byte, header.Size),
		Name: header.Filename,
	}
	_, err := io.ReadFull(file, img.Data)
	if err != nil {
		return nil, pkgErrors.WithMessage(errors.ErrInvalidForm, err.Error())
	}

	if ok := CheckImageContentType(http.DetectContentType(img.Data)); !ok {
		return nil, errors.ErrWrongContentType
	}
	return &img, nil
}
