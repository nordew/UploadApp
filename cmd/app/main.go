package main

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/config"
	v1 "github.com/nordew/UploadApp/internal/controller/http/v1"
	"github.com/nordew/UploadApp/internal/domain/service"
	"github.com/nordew/UploadApp/pkg/auth"
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
	client, err := mongo.NewClient(context.TODO())
	if err != nil {
		logger.Error("failed to connect to db: ", err)
	}

	// Repositories
	userRepository := mongodb.NewUserStorage(client.Database(cfg.DBName).Collection("users"))

	// Service Deps
	hasher := hasher.NewPasswordHasher(cfg.Salt)
	token := auth.NewAuth(jwt.Token{}, cfg.Secret, time.Now().Add(15*time.Minute))

	// Services
	userService := service.NewUserService(userRepository, hasher, token)

	// Handler
	handler := v1.NewHandler(*userService, logger)

	if err := handler.Init(PORT); err != nil {
		logger.Error("failed to init router: ", err)
	}
}
