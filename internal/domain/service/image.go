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
)

const (
	stepOptimization = 3
)

type Images interface {
	Upload(ctx context.Context, images image.Image) (string, error)
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
	imagesRendered, quality := ImageQuality(image)

	generatedId := uuid.NewString()

	for i, v := range imagesRendered {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, v, nil); err != nil {
			return "", err
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

func ImageQuality(img image.Image) ([]image.Image, []int) {
	photoHeight := uint(img.Bounds().Size().Y)
	photoWidth := uint(img.Bounds().Size().X)

	width := []uint{photoWidth, photoWidth - (photoWidth / 4), photoWidth / 2, photoWidth / 4}
	height := []uint{photoHeight, photoHeight - (photoHeight / 4), photoHeight / 2, photoHeight / 4}
	quality := []int{100, 75, 50, 25}

	return reSize(img, width, height), quality
}

func reSize(img image.Image, width, height []uint) (pictures []image.Image) {
	for i := 0; i <= stepOptimization; i++ {
		pictures = append(pictures, resize.Resize(width[i], height[i], img, resize.Lanczos3))
	}

	return pictures
}
