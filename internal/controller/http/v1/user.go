package v1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/controller/http/dto"
	"log"
	"net/http"
)

func (h *Handler) getUser(c *gin.Context) {
	claims := h.getAccessTokenFromRequest(c)

	id := c.Param("sub")

	log.Println(claims.Sub)
	log.Println(id)

	if claims.Sub != id {
		writeErrorResponse(c, http.StatusUnauthorized, "Authentication error", "User ID in the token does not match the requested user ID.")
		return
	}

	user, err := h.userService.GetCredentials(context.Background(), id, false)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "GetCredantials failed", err.Error())
	}

	response := gin.H{
		"user": gin.H{
			"name":            user.Name,
			"email":           user.Email,
			"photos_uploaded": user.PhotosUploaded,
		},
		"subscription": gin.H{},
	}

	writeResponse(c, http.StatusOK, response)
}

func (h *Handler) changePassword(c *gin.Context) {
	var passDto dto.ChangePasswordDTO

	if err := c.ShouldBindJSON(&passDto); err != nil {
		invalidJSONResponse(c)
		return
	}

	err := h.userService.ChangePassword(context.Background(), passDto.Email, passDto.OldPassword, passDto.NewPassword)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, "failed to change password", err.Error())
		return
	}

	writeResponse(c, http.StatusOK, gin.H{})
}
