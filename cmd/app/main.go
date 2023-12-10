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
	"github.com/nordew/UploadApp/pkg/payment"
	"github.com/stripe/stripe-go/v76"
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
	authStorage := psqldb.NewAuthStorage(postgresClient)
	imageStorage := miniodb.NewImageStorage(minioClient, "images")

	hasher := hasher.NewPasswordHasher(cfg.Salt)
	authenticator := auth.NewAuth()

	imageService := service.NewImageService(imageStorage)
	authService := service.NewAuthService(authStorage)
	userService := service.NewUserService(authService, userStorage, hasher, authenticator, cfg.Secret)

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

	stripe.Key = cfg.StripeKey
	stripePayment := payment.NewPayement()

	ctx, cancel := context.WithCancel(context.Background())

	consumer := controller.NewConsumer(channel, q, logger, imageService)

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			logger.Error("failed to consume")
		}
	}()

	handler := v1.NewHandler(userService, imageService, authService, logger, channel, authenticator, stripePayment)
	router := handler.Init()

	go func() {
		if err := server.Run(router, cfg.ServerPort); err != nil {
			logger.Error("failed to run router")
			cancel()
		}
	}()

	select {}
}
