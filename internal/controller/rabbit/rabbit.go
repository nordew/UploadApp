package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nordew/UploadApp/internal/domain/service"
	"github.com/streadway/amqp"
	"image"
	"log/slog"
)

type Consumer struct {
	channel      *amqp.Channel
	queue        amqp.Queue
	logger       *slog.Logger
	imageService service.Images
}

func NewConsumer(channel *amqp.Channel, queue amqp.Queue, logger *slog.Logger, imageService service.Images) *Consumer {
	return &Consumer{
		channel:      channel,
		queue:        queue,
		logger:       logger,
		imageService: imageService,
	}
}

func (c *Consumer) Consume(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		c.logger.Error("Consume() error: ", err)
		return err
	}

	for d := range msgs {
		var imgBytes []byte
		if err := json.Unmarshal(d.Body, &imgBytes); err != nil {
			c.logger.Error("Unmarshal() error: ", err)
			return err
		}

		img, _, err := image.Decode(bytes.NewReader(imgBytes))
		if err != nil {
			c.logger.Error("Decode() error: ", err)
			return err
		}

		_, err = c.imageService.Upload(ctx, img)
		if err != nil {
			c.logger.Error("Upload() error: ", err)
			return err
		}
	}

	return nil
}