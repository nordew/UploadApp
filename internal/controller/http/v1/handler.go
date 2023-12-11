package v1

import (
	"github.com/nordew/UploadApp/pkg/payment"
	"log/slog"
	"net/http"

	"github.com/nordew/UploadApp/pkg/auth"

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
	payment      payment.Payment
}

func NewHandler(userService service.Users, imageService service.Images, logger *slog.Logger, channel *amqp.Channel, auth auth.Authenticator, payment payment.Payment) *Handler {
	return &Handler{
		userService:  userService,
		imageService: imageService,
		logger:       logger,
		channel:      channel,
		auth:         auth,
		payment:      payment,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
	}

	profile := router.Group("profile")
	profile.Use(h.AuthMiddleware())
	{
		profile.GET("/:id", h.getUser)
	}

	image := router.Group("/images")
	image.Use(h.AuthMiddleware())
	{
		image.POST("/upload", h.upload)
		image.GET("/get-all", h.getAll)
		image.GET("/get", h.getBySize)
	}

	payment := router.Group("/payment")
	payment.Use(h.AuthMiddleware())
	{
		payment.GET("/create-intent", h.createPaymentIntent)
	}

	return router
}

func writeResponse(c *gin.Context, statusCode int, h gin.H) {
	c.JSON(statusCode, h)
}

func invalidJSONResponse(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
}

func writeErrorResponse(c *gin.Context, statusCode int, error string, errorDesc string) {
	response := gin.H{
		"error":             error,
		"error_descriptipn": errorDesc,
	}

	c.JSON(statusCode, response)
}

func (h *Handler) getRefreshTokenFromCookie(c *gin.Context) (*auth.ParseTokenClaimsOutput, error) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		writeErrorResponse(c, http.StatusBadRequest, "cookie error", err.Error())
		return nil, err
	}

	claims, err := h.auth.ParseToken(cookie)
	if err != nil {
		writeErrorResponse(c, http.StatusUnauthorized, "Access error", err.Error())
		return nil, err
	}

	return claims, nil
}
