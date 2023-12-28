package common

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {
	config := zap.NewDevelopmentConfig()
	if viper.Get("loglevel") != nil {
		// if it is not set, by default will be 0 - info
		config.Level.SetLevel(zapcore.Level(viper.GetInt("log_level")))
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapLogger, err := config.Build()
	// trunk-ignore(golangci-lint/errcheck)
	defer zapLogger.Sync()
	if err != nil {
		panic(err)
	}
	Logger = zapLogger.Sugar()
}
