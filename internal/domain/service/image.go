package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	miniodb "github.com/nordew/UploadApp/internal/adapters/db/minio"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

const (
	stepOptimization = 3
)

type Images interface {
	Upload(ctx context.Context, images []image.Image) error
}

type ImageService struct {
	storage miniodb.ImageStorage
}

func NewImageService(storage miniodb.ImageStorage) *ImageService {
	return &ImageService{
		storage: storage,
	}
}

func (s *ImageService) Upload(ctx context.Context, images []image.Image) error {
	for _, v := range images {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, v, nil); err != nil {
			return err
		}

		reader := bytes.NewReader(buf.Bytes())

		image := entity.Image{
			ID:     uuid.NewString(),
			Name:   fmt.Sprintf("image"),
			Size:   reader.Size(),
			Reader: reader,
		}

		sizeInMB := float64(image.Size) / (1024 * 1024)
		if sizeInMB > 10 {
			return errors.New("file is bigger than 10mb")
		}

		if err := s.storage.Upload(ctx, image); err != nil {
			return err
		}
	}

	return nil
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
