package main

import (
	"errors"
	"github.com/Rabbit-Go/logger"
	"github.com/Sirupsen/logrus"
	"github.com/evalphobia/logrus_sentry"
	"github.com/spf13/viper"
)

func initConf(n int) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Log.Fatal(err, "No configuration file loaded - using defaults")
	}
	viper.SetDefault("redisServer", "redis://127.0.0.1:6379")
	viper.SetDefault("redisPassword", "")
	viper.SetDefault("mongodb", "mongodb://localhost:27017/keepdb")
	viper.SetDefault("amqpUrl", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("sentryDsn", "")
	viper.SetDefault("workerCount", n)
	if viper.GetString("exchangeName") == "" {
		logger.Log.Fatal(errors.New("ExchangeName Missing"), "Need set exchangename, No default")
	}
	if viper.GetString("sentryDsn") != "" {
		logger.Log.Info("Connect to sentry ", viper.GetString("sentryDsn"))
		hook, err := logrus_sentry.NewSentryHook(viper.GetString("sentryDsn"), []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		hook.StacktraceConfiguration.Enable = true
		hook.StacktraceConfiguration.Context = 5
		if err == nil {
			logger.Log.Hooks.Add(hook)
		} else {
			logger.Log.Fatal(err)
		}
	}
	logger.Log.WithFields(logrus.Fields{"amqpUrl": viper.GetString("amqpUrl")}).Info("AMQP URL")
	logger.Log.WithFields(logrus.Fields{"exchangeName": viper.GetString("exchangeName")}).Info("Exchange Name")
}
