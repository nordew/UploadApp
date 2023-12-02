package v1

import (
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
			writeResponse(c, http.StatusBadRequest, "Invalid or missing token in the Authorization header")
			c.Abort()
			return
		}

		_, err = h.auth.ParseToken(token)
		if err != nil {
			handleTokenValidationError(c, err)
			return
		}
	}
}

func handleTokenValidationError(c *gin.Context, err error) {
	var validationError *jwt.ValidationError
	if errors.As(err, &validationError) {
		switch {
		case validationError.Errors&jwt.ValidationErrorSignatureInvalid != 0:
			writeResponse(c, http.StatusUnauthorized, "Invalid token signature")
		case validationError.Errors&jwt.ValidationErrorExpired != 0:
			writeResponse(c, http.StatusUnauthorized, "Token has expired")
		default:
			writeResponse(c, http.StatusInternalServerError, "Failed to parse token")
		}
	} else {
		writeResponse(c, http.StatusInternalServerError, "Failed to parse token")
	}
}

func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty or missing Authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid Authorization header format")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("empty token")
	}

	return headerParts[1], nil
}
