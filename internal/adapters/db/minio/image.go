package miniodb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

var (
	ErrObjectNotFound = errors.New("object wasn't found")
)

// ImageStorage is the interface that defines methods for storing and retrieving images in an object storage service.
type ImageStorage interface {
	// Upload uploads the provided image to the storage service with the specified identifier.
	// It returns an error if the upload fails.
	Upload(ctx context.Context, image entity.Image) error

	// GetAll retrieves all images associated with the specified identifier from the storage service.
	// It returns a slice of Image entities and an error if the retrieval fails.
	GetAll(ctx context.Context, id string) ([]entity.Image, error)

	// GetBySize retrieves an image of the specified size associated with the given identifier from the storage service.
	// It returns the requested Image entity and an error if the retrieval fails.
	GetBySize(ctx context.Context, id string, size int) (*entity.Image, error)

	// DeleteAllImages deletes all images associated with the specified identifier.
	// The identifier is used to uniquely identify the set of images to delete.
	// It returns an error if the deletion operation fails.
	// The error may indicate issues with object removal or connectivity with the storage service.
	DeleteAllImages(ctx context.Context, id string) error
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

func (s *imageStorage) GetAll(ctx context.Context, id string) ([]entity.Image, error) {
	imageCh := s.db.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Prefix: id,
	})

	var (
		imagesMu sync.Mutex
		images   []entity.Image
		wg       sync.WaitGroup
		errCh    = make(chan error, 1)
	)

	for object := range imageCh {
		if object.Err != nil {
			errCh <- object.Err
			break
		}

		wg.Add(1)

		go func(object minio.ObjectInfo) {
			defer wg.Done()

			objectData, err := s.db.GetObject(ctx, s.bucketName, object.Key, minio.GetObjectOptions{})
			if err != nil {
				errCh <- err
				return
			}
			defer objectData.Close()

			buf := new(bytes.Buffer)
			if _, err := io.Copy(buf, objectData); err != nil {
				errCh <- err
				return
			}

			image := entity.Image{
				ID:     id,
				Name:   object.Key,
				Size:   int64(buf.Len()),
				Reader: bytes.NewReader(buf.Bytes()),
			}

			imagesMu.Lock()
			images = append(images, image)
			imagesMu.Unlock()
		}(object)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return images, nil
}

func (s *imageStorage) GetBySize(ctx context.Context, id string, size int) (*entity.Image, error) {
	prompt := fmt.Sprintf("%s_%d", id, size)

	image, err := s.db.GetObject(ctx, s.bucketName, prompt, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w", ErrObjectNotFound)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, image); err != nil {
		return nil, err
	}

	preparedImage := entity.Image{
		ID:     id,
		Name:   prompt,
		Size:   int64(buf.Len()),
		Reader: bytes.NewReader(buf.Bytes()),
	}

	return &preparedImage, nil
}

func (s *imageStorage) DeleteAllImages(ctx context.Context, id string) error {
	if err := s.db.RemoveObject(ctx, s.bucketName, id, minio.RemoveObjectOptions{}); err != nil {
		return err
	}

	return nil
}
