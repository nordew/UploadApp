package v1

import (
	"bytes"
	"github.com/streadway/amqp"
	"image"
	"image/jpeg"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) upload(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		h.logger.Error("error: %s", err)
		writeResponse(c, http.StatusBadRequest, "invalid file")

		return
	}

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

	var imgBytes []byte
	buffer := bytes.NewBuffer(imgBytes)
	err = jpeg.Encode(buffer, img, nil)
	if err != nil {
		writeResponse(c, http.StatusInternalServerError, "failed to encode image")
		return
	}

	err = h.channel.Publish(
		"",
		"image",
		false,
		false,
		amqp.Publishing{
			ContentType: "image/jpeg",
			Body:        imgBytes,
		})
	if err != nil {
		writeResponse(c, http.StatusInternalServerError, "failed to publish message")
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to queue"})
}
