package main

import (
	"context"
	miniodb "github.com/nordew/UploadApp/internal/adapters/db/minio"
	psqldb "github.com/nordew/UploadApp/internal/adapters/db/postgres"
	"github.com/nordew/UploadApp/internal/config"
	v1 "github.com/nordew/UploadApp/internal/controller/http/v1"
	controller "github.com/nordew/UploadApp/internal/controller/rabbit"
	"github.com/nordew/UploadApp/internal/controller/server"
	"github.com/nordew/UploadApp/internal/domain/service"
	"github.com/nordew/UploadApp/pkg/auth"
	"github.com/nordew/UploadApp/pkg/client/minio"
	"github.com/nordew/UploadApp/pkg/client/psql"
	"github.com/nordew/UploadApp/pkg/client/rabbit"
	"github.com/nordew/UploadApp/pkg/hasher"
	"github.com/nordew/UploadApp/pkg/logging"
)

func main() {
	logger := logging.NewLogger()

	cfg, err := config.NewConfig("main", "yml", "./configs")
	if err != nil {
		logger.Error("failed to create config: ", err)
	}

	postgresClient, err := psql.NewPostgres(&psql.ConnectionInfo{
		Host:     cfg.PGHost,
		Port:     cfg.PGPort,
		User:     cfg.PGUser,
		DBName:   cfg.PGDBName,
		SSLMode:  cfg.PGSSLMode,
		Password: cfg.PGPassword,
	})

	if err != nil {
		logger.Error("failed to connect to postgres: ", err)
		return
	}

	minioClient, err := minio.NewMinioClient(cfg.MinioHost, cfg.MinioUser, cfg.MinioPassword, false, cfg.MinioPort)
	if err != nil {
		logger.Error("failed to connect to minio: ", err)
		return
	}

	userStorage := psqldb.NewUserStorage(postgresClient)
	imageStorage := miniodb.NewImageStorage(minioClient, "images", logger)
	dashboardStorage := psqldb.NewDashboardStorage(postgresClient, logger)

	hasher := hasher.NewPasswordHasher(cfg.Salt)
	authenticator := auth.NewAuth(logger)

	imageService := service.NewImageService(imageStorage, logger)
	userService := service.NewUserService(userStorage, hasher, authenticator, logger, cfg.Secret)
	dashboardService := service.NewDashboardService(dashboardStorage)

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
		logger.Error("failed to declare queue")
	}

	ctx, cancel := context.WithCancel(context.Background())

	consumer := controller.NewConsumer(channel, q, logger, imageService, dashboardService, userService)

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			logger.Error("failed to consume")
		}
	}()

	handler := v1.NewHandler(userService, imageService, dashboardService, logger, channel, authenticator)
	router := handler.Init()

	go func() {
		if err := server.Run(router, cfg.ServerPort); err != nil {
			logger.Error("failed to run router")
			cancel()
		}
	}()

	select {}
}
