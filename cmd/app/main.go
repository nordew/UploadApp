package main

import (
	"context"
	miniodb "github.com/nordew/UploadApp/internal/adapters/db/minio"
	psqldb "github.com/nordew/UploadApp/internal/adapters/db/postgres"
	"github.com/nordew/UploadApp/internal/config"
	v1 "github.com/nordew/UploadApp/internal/controller/http/v1"
	controller "github.com/nordew/UploadApp/internal/controller/rabbit"
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

	// Storages
	userStorage := psqldb.NewUserStorage(postgresClient)
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

	// Stripe
	stripe.Key = cfg.StripeKey
	stripePayment := payment.NewPayement()

	// Create a context that can be canceled to signal shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals to gracefully shutdown
	//go func() {
	//	sigCh := make(chan os.Signal, 1)
	//	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	//	<-sigCh
	//	logger.Info("Received shutdown signal")
	//	cancel()
	//}()

	consumer := controller.NewConsumer(channel, q, logger, imageService)

	go func() {
		if err := consumer.Consume(ctx); err != nil {
			logger.Error("failed to consume")
		}
	}()

	// Handler
	go func() {
		handler := v1.NewHandler(userService, imageService, logger, channel, authenticator, stripePayment)

		if err := handler.Init(PORT); err != nil {
			logger.Error("failed to init router: ", err)
			cancel()
		}
	}()

	select {}
	//
	//<-ctx.Done()
	//
	//logger.Info("Shutting down...")
	//
	//// Close RabbitMQ connection
	//if err := conn.Close(); err != nil {
	//	logger.Error("failed to close RabbitMQ connection: ", err)
	//}
	//
	//logger.Info("Shutdown complete.")
}
