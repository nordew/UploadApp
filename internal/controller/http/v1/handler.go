package v1

import (
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/nordew/UploadApp/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/nordew/UploadApp/internal/domain/service"
	"github.com/streadway/amqp"
)

type Handler struct {
	imageService     service.Images
	userService      service.Users
	dashboardService service.Dashboards
	logger           *logrus.Logger
	channel          *amqp.Channel
	auth             auth.Authenticator
}

func NewHandler(
	userService service.Users,
	imageService service.Images,
	dashboardService service.Dashboards,
	logger *logrus.Logger,
	channel *amqp.Channel,
	auth auth.Authenticator) *Handler {
	return &Handler{
		userService:      userService,
		imageService:     imageService,
		dashboardService: dashboardService,
		logger:           logger,
		channel:          channel,
		auth:             auth,
	}
}

func (h *Handler) Init() *gin.Engine {
	router := gin.Default()

	root := router.Group("/")
	{
		root.POST("/change-password", h.changePassword)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.GET("/sign-in", h.signIn)
		auth.GET("/refresh", h.refresh)
	}

	profile := router.Group("/profile")
	profile.Use(h.AuthMiddleware())
	{
		profile.GET("/:sub", h.getUser)
	}

	image := router.Group("/images")
	image.Use(h.AuthMiddleware())
	{
		image.POST("/upload", h.upload)
		image.GET("/all", h.getAllImages)
		image.GET("/by-size", h.getBySize)
		image.DELETE("/delete/:name", h.deleteAllImages)
	}

	dashboard := router.Group("/dashboard")
	dashboard.Use(h.AuthAdminMiddleware())
	{
		dashboard.GET("/logs", h.getLogs)
		dashboard.DELETE("/logs/:id", h.deleteLog)
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
		"error_description": errorDesc,
	}

	c.JSON(statusCode, response)
}

func (h *Handler) getAccessTokenFromRequest(c *gin.Context) *auth.ParseTokenClaimsOutput {
	logger := h.logger.WithField("function", "getAccessTokenFromRequest")

	accessToken := extractTokenFromHeader(c.Request.Header, "Authorization")

	if accessToken == "" {
		logger.Error("access token not provided in headers")
		writeErrorResponse(c, http.StatusUnauthorized, "auth", "access token not provided in headers")
		c.Abort()
		return nil
	}

	accessTokenClaims, err := h.auth.ParseToken(accessToken)
	if err != nil {
		logger.WithError(err).Error("failed to parse access token")
		writeErrorResponse(c, http.StatusUnauthorized, "auth", "failed to parse access token")
		c.Abort()
		return nil
	}

	return accessTokenClaims
}

func (h *Handler) getRefreshTokenFromRequest(c *gin.Context) *auth.ParseTokenClaimsOutput {
	logger := h.logger.WithField("function", "getRefreshTokenFromRequest")

	refreshToken := extractTokenFromHeader(c.Request.Header, "Refresh-Token")

	if refreshToken == "" {
		logger.Error("refresh token not provided in headers")
		writeErrorResponse(c, http.StatusUnauthorized, "auth", "refresh token not provided in headers")
		c.Abort()
		return nil
	}

	refreshTokenClaims, err := h.auth.ParseToken(refreshToken)
	if err != nil {
		logger.WithError(err).Error("failed to parse refresh token")
		writeErrorResponse(c, http.StatusUnauthorized, "auth", "failed to parse refresh token")
		c.Abort()
		return nil
	}

	return refreshTokenClaims
}

func extractTokenFromHeader(headers map[string][]string, headerKey string) string {
	token := ""
	if headerValues, exists := headers[headerKey]; exists && len(headerValues) > 0 {
		token = headerValues[0]
	}
	return token
}
