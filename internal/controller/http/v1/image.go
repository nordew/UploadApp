package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"sync"

	"github.com/nordew/UploadApp/internal/controller/http/dto"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/streadway/amqp"

	"github.com/gin-gonic/gin"
)

func (h *Handler) uploadImage(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "Failed to parse form")
		return
	}

	files, ok := c.Request.MultipartForm.File["photo"]
	if !ok || len(files) == 0 {
		writeErrorResponse(c, http.StatusBadRequest, "image", "No file found in the 'photo' field")
		return
	}

	for _, file := range files {
		if err := h.processFile(c, file); err != nil {
			return
		}
	}

	response := gin.H{
		"success": "adeed to queue",
	}

	writeResponse(c, http.StatusCreated, response)
}

func (h *Handler) processFile(c *gin.Context, file *multipart.FileHeader) error {
	openedFile, err := file.Open()
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "Failed to open file")
		return err
	}
	defer openedFile.Close()

	content, err := io.ReadAll(openedFile)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "Failed to read file")
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(content))
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "Failed to convert image")
		return err
	}

	var imgBytesBuffer bytes.Buffer
	if err := jpeg.Encode(&imgBytesBuffer, img, nil); err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "Failed to convert image to bytes")
		return err
	}

	imgBytes := imgBytesBuffer.Bytes()

	if err := h.publishImageToQueue(c, imgBytes); err != nil {
		return err
	}

	return nil
}

func (h *Handler) publishImageToQueue(c *gin.Context, imgBytes []byte) error {
	marshalledImg, err := json.Marshal(&imgBytes)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "Failed to marshall image")
		return err
	}

	err = h.channel.Publish(
		"",
		"image",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        marshalledImg,
		})
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "Failed to add image to queue")
		return err
	}

	return nil
}

func (h *Handler) getAllImages(c *gin.Context) {
	var getAllImageDTO dto.GetAllImageDTO

	if err := c.ShouldBindJSON(&getAllImageDTO); err != nil {
		invalidJSONResponse(c)
		return
	}

	images, err := h.imageService.GetAll(context.Background(), getAllImageDTO.ID)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to get photos", err.Error())
		return
	}

	c.Writer.Header().Set("Content-Type", "multipart/x-mixed-replace; boundary=frame")

	var wg sync.WaitGroup

	for _, v := range images {
		wg.Add(1)

		go func(entityImage entity.Image) {
			defer wg.Done()

			img, _, err := image.Decode(entityImage.Reader)
			if err != nil {
				writeErrorResponse(c, http.StatusInternalServerError, "image", "failed to decode image")
				return
			}

			var encodedImageBuffer bytes.Buffer
			if err := jpeg.Encode(&encodedImageBuffer, img, nil); err != nil {
				writeErrorResponse(c, http.StatusInternalServerError, "image", "failed to encode image")
				return
			}

			c.Data(http.StatusOK, "image/jpeg", encodedImageBuffer.Bytes())
		}(v)
	}

	wg.Wait()

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) deleteAllImages(c *gin.Context) {
	name := c.Param("name")

	if err := h.imageService.DeleteAllImages(context.Background(), name); err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to delete image", err.Error())
	}

	writeResponse(c, http.StatusOK, gin.H{})
}

func (h *Handler) getBySize(c *gin.Context) {
	var GetImageBySizeDTO dto.GetImageBySizeDTO

	if err := c.ShouldBindJSON(&GetImageBySizeDTO); err != nil {
		invalidJSONResponse(c)
		return
	}

	entityImg, err := h.imageService.GetBySize(context.Background(), GetImageBySizeDTO.ID, GetImageBySizeDTO.Size)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image (failed to make getBySize())", err.Error())
		return
	}

	img, _, err := image.Decode(entityImg.Reader)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "failed to decode image")
		return
	}

	var encodedImageBuffer bytes.Buffer
	if err := jpeg.Encode(&encodedImageBuffer, img, nil); err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "image", "failed to encode image")
		return
	}

	c.Data(http.StatusOK, "image/jpeg", encodedImageBuffer.Bytes())
}
