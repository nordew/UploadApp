package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.SetOutput(os.Stdout)

	return logger
}
