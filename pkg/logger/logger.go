package logger

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func New() logr.Logger {
	opts := zap.Options{
		Development: true,
		Level:       zapcore.PanicLevel,
	}

	if viper.GetBool("debug") {
		opts.Level = zapcore.DebugLevel
	}

	if viper.GetBool("verbose") {
		opts.Level = zapcore.InfoLevel
	}

	logger := zap.New(zap.UseFlagOptions(&opts))
	logger.WithName("rasaxctl")

	return logger
}
