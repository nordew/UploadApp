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

type Images interface {
	Upload(ctx context.Context, images image.Image) (string, error)
	GetAll(ctx context.Context, id string) ([]entity.Image, error)
	GetBySize(ctx context.Context, id string, size int) (*entity.Image, error)
}

type ImageService struct {
	storage miniodb.ImageStorage
}

func NewImageService(storage miniodb.ImageStorage) *ImageService {
	return &ImageService{
		storage: storage,
	}
}

func (s *ImageService) Upload(ctx context.Context, image image.Image) (string, error) {
	imagesRendered, quality, err := ImageQuality(image)
	if err != nil {
		return "", err
	}

	generatedId := uuid.NewString()

	for i, v := range imagesRendered {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, v, nil); err != nil {
			return "", fmt.Errorf("failed to encode image")
		}

		reader := bytes.NewReader(buf.Bytes())

		idFormatted := fmt.Sprintf("%s_%d.jpg", generatedId, quality[i])

		resImage := entity.Image{
			Name:   idFormatted,
			Size:   reader.Size(),
			Reader: reader,
		}

		if err := s.storage.Upload(ctx, resImage); err != nil {
			log.Printf("upload error: %s", err)
			return "", err
		}
	}

	return generatedId, nil
}

func (s *ImageService) GetAll(ctx context.Context, id string) ([]entity.Image, error) {
	return s.storage.GetAll(ctx, id)
}

func (s *ImageService) GetBySize(ctx context.Context, id string, size int) (*entity.Image, error) {
	return s.storage.GetBySize(ctx, id, size)
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
