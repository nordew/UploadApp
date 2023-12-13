package service

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	miniodb "github.com/nordew/UploadApp/internal/adapters/db/minio"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/pkg/errors"
)

const (
	stepOptimization = 3
)

// Images is the interface that defines methods for interacting with image-related operations.
// It provides functionalities for uploading, retrieving, and manipulating images.
type Images interface {
	// Upload uploads the given image and returns a unique identifier for the uploaded image.
	// The image is processed to create multiple versions with different qualities.
	// It returns the generated identifier and an error if the upload fails.
	Upload(ctx context.Context, image image.Image, userId string) error

	// GetAll retrieves all images associated with the specified identifier.
	// It returns a slice of Image entities and an error if the retrieval fails.
	GetAll(ctx context.Context, id string) ([]entity.Image, error)

	// GetBySize retrieves an image of the specified size associated with the given identifier.
	// It returns the requested Image entity and an error if the retrieval fails.
	GetBySize(ctx context.Context, id string, size int) (*entity.Image, error)

	// DeleteAllImages deletes all versions of images associated with the specified identifier.
	// It removes images with different qualities (e.g., different compression levels) to clean up storage.
	// The method takes a context and an identifier as input and returns an error if the deletion fails.
	DeleteAllImages(ctx context.Context, id string) error
}

type ImageService struct {
	storage miniodb.ImageStorage
}

func NewImageService(storage miniodb.ImageStorage) *ImageService {
	return &ImageService{
		storage: storage,
	}
}

func (s *ImageService) Upload(ctx context.Context, image image.Image, userId string) error {
	imagesRendered, quality, err := ImageQuality(image)
	if err != nil {
		return err
	}

	generatedId := uuid.NewString()

	for i, v := range imagesRendered {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, v, nil); err != nil {
			return fmt.Errorf("failed to encode image")
		}

		reader := bytes.NewReader(buf.Bytes())

		idFormatted := fmt.Sprintf("%s_%d_%s.jpeg", generatedId, quality[i], userId)

		resImage := entity.Image{
			Name:   idFormatted,
			Size:   reader.Size(),
			Reader: reader,
		}

		if err := s.storage.Upload(ctx, resImage); err != nil {
			log.Printf("upload error: %s", err)
			return err
		}
	}

	return nil
}

func (s *ImageService) GetAll(ctx context.Context, id string) ([]entity.Image, error) {
	return s.storage.GetAll(ctx, id)
}

func (s *ImageService) GetBySize(ctx context.Context, id string, size int) (*entity.Image, error) {
	return s.storage.GetBySize(ctx, id, size)
}

func (s *ImageService) DeleteAllImages(ctx context.Context, id string) error {
	quality := []string{"100", "75", "50", "25"}

	for _, v := range quality {
		if err := s.storage.DeleteAllImages(ctx, id, v); err != nil {
			return err
		}
	}

	return nil
}

func ImageQuality(img image.Image) ([]image.Image, []int, error) {
	if img == nil {
		return nil, nil, errors.New("input image is nil")
	}

	photoHeight := uint(img.Bounds().Size().Y)
	photoWidth := uint(img.Bounds().Size().X)

	if photoHeight == 0 || photoWidth == 0 {
		return nil, nil, errors.New("invalid image dimensions")
	}

	width := []uint{photoWidth, photoWidth - (photoWidth / 4), photoWidth / 2, photoWidth / 4}
	height := []uint{photoHeight, photoHeight - (photoHeight / 4), photoHeight / 2, photoHeight / 4}
	quality := []int{100, 75, 50, 25}

	resizedImages, err := reSize(img, width, height)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resize images: %w", err)
	}

	return resizedImages, quality, nil
}

func reSize(img image.Image, width, height []uint) ([]image.Image, error) {
	if img == nil {
		return nil, errors.New("input image is nil")
	}

	var pictures []image.Image
	for i := 0; i <= stepOptimization; i++ {
		resizedImg := resize.Resize(width[i], height[i], img, resize.Lanczos3)
		if resizedImg == nil {
			return nil, fmt.Errorf("failed to resize image at index %d", i)
		}
		pictures = append(pictures, resizedImg)
	}

	return pictures, nil
}
