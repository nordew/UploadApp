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

	id, err := h.imageService.Upload(context.TODO(), img)
	if err != nil {
		writeResponse(c, http.StatusInternalServerError, "failed to upload image")

		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
