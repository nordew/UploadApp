package v1

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/domain/service"
)

type Handler struct {
	userService service.UserService
	logger      *slog.Logger
}

func NewHandler(userService service.UserService, logger *slog.Logger) *Handler {
	return &Handler{
		userService: userService,
		logger:      logger,
	}
}

func (h *Handler) Init(port string) error {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
	}

	if err := router.Run(port); err != nil {
		return err
	}

	return nil
}

func writeResponse(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{"message": msg})
}
