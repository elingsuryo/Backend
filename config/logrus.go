package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewLogrus creates a new logrus logger
func NewLogrus(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	env := viper.GetString("APP_ENV")

	if env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}

	logLevel := viper.GetInt("APP_LOG_LEVEL")
	log.SetLevel(logrus.Level(logLevel))

	return log
}
