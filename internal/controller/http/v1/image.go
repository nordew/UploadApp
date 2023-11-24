package v1

import (
	"bytes"
	"encoding/json"
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

	c.JSON(http.StatusOK, gin.H{"message": "added to queue"})
}
