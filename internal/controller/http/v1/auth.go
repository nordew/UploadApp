package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/domain/entity"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input entity.SignUpInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.WithError(err).Error("signUp: invalid JSON body")
		invalidJSONResponse(c)
		return
	}

	if err := h.userService.SignUp(context.TODO(), input); err != nil {
		h.logger.WithError(err).Error("signUp: failed to SignUp")
		writeErrorResponse(c, http.StatusInternalServerError, "failed to SignUp", err.Error())
		return
	}

	writeResponse(c, http.StatusCreated, gin.H{})
}

func (h *Handler) signIn(c *gin.Context) {
	var input entity.SignInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.WithError(err).Error("signIn: failed to parse data")
		invalidJSONResponse(c)
		return
	}

	accessToken, refreshToken, err := h.userService.SignIn(context.Background(), input)
	if err != nil {
		h.logger.WithError(err).Error("signIn: failed to SignIn")
		writeErrorResponse(c, http.StatusInternalServerError, "failed to SignIn", err.Error())
		return
	}

	writeTokensInHeaders(c, accessToken, refreshToken)
}

func (h *Handler) refresh(c *gin.Context) {
	claims := h.getRefreshTokenFromRequest(c)

	accessToken, refreshToken, err := h.userService.Refresh(context.Background(), claims.Sub, claims.Role)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to generate new tokens", err.Error())
		return
	}

	writeTokensInHeaders(c, accessToken, refreshToken)
}

func writeTokensInHeaders(c *gin.Context, accessToken, refreshToken string) {
	c.Header("Access-Token", accessToken)
	c.Header("Refresh-Token", refreshToken)

	writeResponse(c, http.StatusOK, gin.H{})
}
