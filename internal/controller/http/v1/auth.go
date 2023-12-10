package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		writeErrorResponse(c, http.StatusBadRequest, "cookie error", err.Error())
		return
	}

	accessToken, refreshToken, err := h.authService.RefreshTokens(context.Background(), cookie)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to generate new tokens", err.Error())
		return
	}

	c.SetCookie("access_token", accessToken, 0, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 0, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{})
}
