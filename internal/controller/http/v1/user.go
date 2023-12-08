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
		invalidJSONResponse(c)
		h.logger.Error("signUp: invalid JSON body ", err)
		return
	}

	id, err := h.userService.SignUp(context.TODO(), input)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to SignUp", err.Error())
		h.logger.Debug(err.Error())
		return
	}

	response := gin.H{
		"id": id,
	}

	writeResponse(c, http.StatusCreated, response)
	h.logger.Debug("signUp: user was created ")
}

func (h *Handler) signIn(c *gin.Context) {
	var input entity.SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		invalidJSONResponse(c)
		h.logger.Error("signUp: failed to parse data ", err)
		return
	}

	acccesToken, err := h.userService.SignIn(context.TODO(), input)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to SignIn", err.Error())
		h.logger.Error(err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": acccesToken})
}
