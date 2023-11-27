package v1

import (
	"github.com/nordew/UploadApp/pkg/auth"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/domain/service"
	"github.com/streadway/amqp"
)

type Handler struct {
	imageService service.Images
	userService  service.Users
	logger       *slog.Logger
	channel      *amqp.Channel
	auth         auth.Authenticator
}

func NewHandler(userService service.Users, imageService service.Images, logger *slog.Logger, channel *amqp.Channel, auth auth.Authenticator) *Handler {
	return &Handler{
		userService:  userService,
		imageService: imageService,
		logger:       logger,
		channel:      channel,
		auth:         auth,
	}
}

func (h *Handler) Init(port string) error {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
	}

	image := router.Group("/images")
	image.Use(h.AuthMiddleware())
	{
		image.POST("/upload", h.upload)
		image.GET("/get-all", h.getAll)
	}

	if err := router.Run(port); err != nil {
		return err
	}

	return nil
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	writeResponse(c, statusCode, message)
	c.Abort()
}

func writeResponse(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{"message": msg})
}
