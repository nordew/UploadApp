package miniodb

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

type ImageStorage interface {
	Upload(ctx context.Context, images entity.Image) error
}

type imageStorage struct {
	db         *minio.Client
	bucketName string
}

func NewImageStorage(db *minio.Client, bucketName string) *imageStorage {
	return &imageStorage{
		db:         db,
		bucketName: bucketName,
	}
}

func (s *imageStorage) Upload(ctx context.Context, image entity.Image) error {
	_, err := s.db.PutObject(ctx, s.bucketName, image.Name, image.Reader, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return err
	}

	return nil
}
