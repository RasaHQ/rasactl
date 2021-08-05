package logger

import (
	"github.com/RasaHQ/rasaxctl/pkg/types"
	"github.com/go-logr/logr"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func New(flags *types.RasaXCtlFlags) logr.Logger {
	opts := zap.Options{
		Development: true,
		Level:       zapcore.PanicLevel,
	}

	if flags.Global.Debug {
		opts.Level = zapcore.DebugLevel
	}

	if flags.Global.Verbose {
		opts.Level = zapcore.InfoLevel
	}

	logger := zap.New(zap.UseFlagOptions(&opts))
	logger.WithName("rasaxctl")

	return logger
}
