package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/domain/entity"
)

func (h *Handler) signUp(c *gin.Context) {
	var input entity.SignUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		writeResponse(c, http.StatusBadRequest, "invalid JSON body")
		h.logger.Error("signUp: invalid JSON body ", err)
		return
	}

	if err := h.userService.SignUp(context.TODO(), input); err != nil {
		writeResponse(c, http.StatusInternalServerError, err.Error())
		h.logger.Debug(err.Error())
		return
	}

	writeResponse(c, http.StatusCreated, fmt.Sprintf("user %s was created", input.Name))
	h.logger.Debug("signUp: user was created ")
}

func (h *Handler) signIn(c *gin.Context) {
	var input entity.SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		writeResponse(c, http.StatusBadRequest, "failed to parse data")
		h.logger.Error("signUp: failed to parse data ", err)
		return
	}

	acccesToken, err := h.userService.SignIn(context.TODO(), input)
	if err != nil {
		writeResponse(c, http.StatusInternalServerError, err.Error())
		h.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": acccesToken})
}
