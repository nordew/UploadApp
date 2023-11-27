package miniodb

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"io"
)

type ImageStorage interface {
	Upload(ctx context.Context, images entity.Image) error
	GetAll(ctx context.Context, id string) ([]entity.Image, error)
	GetBySize(ctx context.Context, id string, size int) (*entity.Image, error)
}

type imageStorage struct {
	db         *minio.Client
	bucketName string
}

func (s *imageStorage) GetBySize(ctx context.Context, id string, size int) (*entity.Image, error) {
	return nil, nil
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

func (s *imageStorage) GetAll(ctx context.Context, id string) ([]entity.Image, error) {
	imageCh := s.db.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: id,
	})

	var images []entity.Image

	for object := range imageCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objectData, err := s.db.GetObject(ctx, s.bucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		defer objectData.Close()

		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, objectData); err != nil {
			return nil, err
		}

		image := entity.Image{
			ID:     id,
			Name:   object.Key,
			Size:   int64(buf.Len()),
			Reader: bytes.NewReader(buf.Bytes()),
		}

		images = append(images, image)
	}

	return images, nil
}
