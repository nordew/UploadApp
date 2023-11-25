package v1

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c)
		if err != nil {
			writeResponse(c, http.StatusBadRequest, err.Error())

			c.Abort()

			return
		}

		claims, err := h.userService.ParseToken(context.Background(), token)
		if err != nil {
			var validationError *jwt.ValidationError
			if errors.As(err, &validationError) {
				switch {
				case validationError.Errors&jwt.ValidationErrorSignatureInvalid != 0:
					handleError(c, http.StatusUnauthorized, "Invalid token signature", err)
				case validationError.Errors&jwt.ValidationErrorExpired != 0:
					handleError(c, http.StatusUnauthorized, "Token has expired", err)
				default:
					handleError(c, http.StatusInternalServerError, "Failed to parse token", err)
				}
			} else {
				handleError(c, http.StatusInternalServerError, "Failed to parse token", err)
			}
			return
		}

		userId := claims["sub"]

		writeResponse(c, http.StatusOK, userId.(string))
	}
}

func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
