package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			writeErrorResponse(c, http.StatusUnauthorized, "Middleware", "no access token")
			c.Abort()
			return
		}

		_, err = h.auth.ParseToken(cookie)
		if err != nil {
			handleTokenValidationError(c, err)
			c.Abort()
			return
		}
	}
}

func handleTokenValidationError(c *gin.Context, err error) {
	var validationError *jwt.ValidationError
	if errors.As(err, &validationError) {
		switch {
		case validationError.Errors&jwt.ValidationErrorSignatureInvalid != 0:
			writeErrorResponse(c, http.StatusUnauthorized, "JWT", "Invalid token signature")
		case validationError.Errors&jwt.ValidationErrorExpired != 0:
			writeErrorResponse(c, http.StatusUnauthorized, "JWT", "Token has expired")
		default:
			writeErrorResponse(c, http.StatusInternalServerError, "JWT", "Failed to parse token")
		}
	} else {
		writeErrorResponse(c, http.StatusInternalServerError, "JWT", "Failed to parse token")
	}
}
