package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

func (h *Handler) signUp(c *gin.Context) {
	var input entity.SignUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		writeResponse(c, http.StatusBadRequest, "failed to parse user")
		h.logger.Error("signUp: failed to parse user ", err)
		return
	}

	if err := h.userService.SignUp(context.TODO(), input); err != nil {
		writeResponse(c, http.StatusInternalServerError, "failed to create user")
		h.logger.Error("signUp: failed to create user ", err)
		return
	}

	writeResponse(c, http.StatusCreated, "user created")
	h.logger.Debug("signUp: user was created ")
}
