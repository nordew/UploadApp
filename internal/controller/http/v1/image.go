package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nordew/UploadApp/internal/controller/http/dto"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"github.com/streadway/amqp"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

func (h *Handler) upload(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		h.logger.Error("Error parsing multipart form: %s", err)
		writeResponse(c, http.StatusInternalServerError, "Failed to parse form")
		return
	}

	files, ok := c.Request.MultipartForm.File["photo"]
	if !ok || len(files) == 0 {
		h.logger.Error("No file found in the 'photo' field")
		writeResponse(c, http.StatusBadRequest, "No file found in the 'photo' field")
		return
	}

	for _, file := range files {
		openedFile, err := file.Open()
		if err != nil {
			writeResponse(c, http.StatusInternalServerError, "failed to open file")

			return
		}
		defer openedFile.Close()

		content, err := io.ReadAll(openedFile)
		if err != nil {
			writeResponse(c, http.StatusInternalServerError, "failed to read file")

			return
		}

		img, _, err := image.Decode(bytes.NewReader(content))
		if err != nil {
			writeResponse(c, http.StatusInternalServerError, "failed to convert image")

			return
		}

		var imgBytesBuffer bytes.Buffer
		if err := jpeg.Encode(&imgBytesBuffer, img, nil); err != nil {
			writeResponse(c, http.StatusInternalServerError, "failed to convert image to bytes")

			return
		}

		imgBytes := imgBytesBuffer.Bytes()

		marshalledImg, err := json.Marshal(&imgBytes)
		if err != nil {
			writeResponse(c, http.StatusInternalServerError, "failed to marshall image")

			return
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
			writeResponse(c, http.StatusInternalServerError, "failed to add image to queue")

			return
		}
	}

	writeResponse(c, http.StatusCreated, "added to queue")
}

func (h *Handler) getAll(c *gin.Context) {
	var getAllImageDTO dto.GetAllImageDTO

	if err := c.ShouldBindJSON(&getAllImageDTO); err != nil {
		writeResponse(c, http.StatusBadRequest, "invalid JSON")

		return
	}

	images, err := h.imageService.GetAll(context.Background(), getAllImageDTO.ID)
	if err != nil {
		writeResponse(c, http.StatusInternalServerError, "")

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
				writeResponse(c, http.StatusInternalServerError, "failed to decode image")

				return
			}

			if err := jpeg.Encode(c.Writer, img, nil); err != nil {
				writeResponse(c, http.StatusInternalServerError, "failed to encode image")

				return
			}

			c.Writer.(http.Flusher).Flush()
		}(v)
	}

	wg.Wait()

	c.AbortWithStatus(http.StatusOK)
}
