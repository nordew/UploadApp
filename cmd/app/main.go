package main

import (
	"context"
	miniodb "github.com/nordew/UploadApp/internal/adapters/db/minio"
	v1 "github.com/nordew/UploadApp/internal/controller/http/v1"
	controller "github.com/nordew/UploadApp/internal/controller/rabbit"
	"github.com/nordew/UploadApp/pkg/client/rabbit"

	"github.com/nordew/UploadApp/internal/adapters/db/mongodb"
	"github.com/nordew/UploadApp/internal/config"
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
		logger.Error("failed to connect to minio: ", err)
	}

	// Storages
	userStorage := mongodb.NewUserStorage(mongoClient.Database(cfg.MongoDBName).Collection("users"))
	imageStorage := miniodb.NewImageStorage(minioClient, "images")

	// Pkg
	hasher := hasher.NewPasswordHasher(cfg.Salt)
	authenticator := auth.NewAuth()

	// Services
	userService := service.NewUserService(userStorage, hasher, authenticator, cfg.Secret)
	imageService := service.NewImageService(imageStorage)

	// Queue
	conn, err := rabbit.NewRabbitClient(cfg.Rabbit)
	if err != nil {
		logger.Error("failed to connect to rabbit: ", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		logger.Error("failed to open channel: ", err)

	}

	q, err := channel.QueueDeclare(
		"image", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logger.Error("failed to declare a queue")
	}

	// Consumer
	consumer := controller.NewConsumer(channel, q, logger, imageService)

	go func() {
		if err := consumer.Consume(context.Background()); err != nil {
			logger.Error("failed to consume")
		}
	}()

	// Handler
	go func() {
		handler := v1.NewHandler(userService, imageService, logger, channel, authenticator)

		if err := handler.Init(PORT); err != nil {
			logger.Error("failed to init router: ", err)
		}
	}()

	select {}
}
