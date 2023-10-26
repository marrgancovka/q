package file_storage

import "hack/internal/models"

type RepoI interface {
	Get(bName, fName string) (*models.S3File, error)
	Upload(file *models.S3File) error
}

type UseCaseI interface {
	Get(bName, fName string) (*models.S3File, error)
	Upload(file *models.S3File) error
}
