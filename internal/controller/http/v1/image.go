package v1

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"image"
	"io"
	"net/http"
)

func (h *Handler) upload(c *gin.Context) {
	var images []image.Image

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

	images = append(images, img)

	h.imageService.Upload(context.TODO(), images)

	writeResponse(c, http.StatusOK, "image posted")
}
