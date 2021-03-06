/*
Copyright © 2021 Rasa Technologies GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package logger

import (
	"github.com/go-logr/logr"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/RasaHQ/rasactl/pkg/types"
)

// New initializes and return a new logger object.
func New(flags *types.RasaCtlFlags) logr.Logger {
	opts := zap.Options{
		Development: true,
		Level:       zapcore.PanicLevel,
	}

	if flags.Global.Debug {
		opts.Level = zapcore.DebugLevel
		opts.TimeEncoder = zapcore.ISO8601TimeEncoder
	}

	if flags.Global.Verbose {
		opts.Level = zapcore.InfoLevel
		opts.TimeEncoder = zapcore.ISO8601TimeEncoder
	}

	logger := zap.New(zap.UseFlagOptions(&opts))
	logger.WithName("rasactl")

	return logger
}
