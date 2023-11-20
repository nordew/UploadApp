package main

import (
	"context"
	"time"

	miniodb "github.com/nordew/UploadApp/internal/adapters/db/minio"

	"github.com/golang-jwt/jwt"
	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/config"
	v1 "github.com/nordew/UploadApp/internal/controller/http/v1"
	"github.com/nordew/UploadApp/internal/domain/service"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/client/minio"
	mongo "github.com/nordew/UploadApp/pkg/client/mongodb"
	"github.com/nordew/UploadApp/pkg/hasher"
	"github.com/nordew/UploadApp/pkg/logging"
)

const (
	PORT = ":8080"
)

func main() {
	// Logging
	logger := logging.NewLogger()

	// Config
	cfg, err := config.NewConfig("main", "yml", "./configs")
	if err != nil {
		logger.Error("failed to create config: ", err)
	}

	// DB
	mongoClient, err := mongo.NewClient(context.TODO())
	if err != nil {
		logger.Error("failed to connect to monggo: ", err)
	}

	minioClient, err := minio.NewMinioClient(cfg.MinioHost, cfg.MinioUser, cfg.MinioPassword, false, cfg.MinioPort)
	if err != nil {
		logger.Error("failed to connect to minio: %s", err)
	}

	// Storages
	userStorage := mongodb.NewUserStorage(mongoClient.Database(cfg.MongoDBName).Collection("users"))
	imageStorage := miniodb.NewImageStorage(minioClient, "images")

	// Service Deps
	hasher := hasher.NewPasswordHasher(cfg.Salt)
	token := auth.NewAuth(jwt.Token{}, cfg.Secret, time.Now().Add(15*time.Minute))

	// Services
	userService := service.NewUserService(userStorage, hasher, token)
	imageService := service.NewImageService(imageStorage)

	// Handler
	handler := v1.NewHandler(userService, imageService, logger)

	if err := handler.Init(PORT); err != nil {
		logger.Error("failed to init router: ", err)
	}
}
