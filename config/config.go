package config

import (
	"jatis_mobile_api/logs"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	RabbitMQURL string
	PostgresURL string
	PORT        int
}

var logger = logs.SetupLogger()

func LoadConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to read configuration file", struct{ Error error }{Error: err})
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		logs.LogWithFields(logger, logrus.ErrorLevel, "Failed to unmarshal configuration", struct{ Error error }{Error: err})
		return config, err
	}

	logs.LogWithFields(logger, logrus.InfoLevel, "Configuration loaded successfully", struct {
		Port        int
		PostgresURL string
		RabbitMQURL string
	}{
		Port:        config.PORT,
		PostgresURL: config.PostgresURL,
		RabbitMQURL: config.RabbitMQURL,
	})

	return config, nil
}
