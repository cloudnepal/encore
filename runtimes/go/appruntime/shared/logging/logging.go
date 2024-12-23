//go:build encore_app

package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"

	"encore.dev/appruntime/exported/config"
	"encore.dev/appruntime/shared/appconf"
	"encore.dev/appruntime/shared/cloud"
)

var RootLogger = configure(appconf.Static, appconf.Runtime)

func configure(static *config.Static, runtime *config.Runtime) zerolog.Logger {
	var logOutput io.Writer = os.Stderr
	if static.PrettyPrintLogs {
		logOutput = zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
			w.Out = logOutput
		})
	}

	level := zerolog.TraceLevel
	if runtime.LogConfig != "" {
		if l, err := zerolog.ParseLevel(runtime.LogConfig); err == nil {
			level = l
		}
	}

	reconfigureZerologFormat(runtime)
	return zerolog.New(logOutput).Level(level).With().Timestamp().Logger()
}

func reconfigureZerologFormat(runtime *config.Runtime) {
	// Note: if updating this function, also update
	// mapCloudFieldNamesToExpected in cli/cmd/encore/logs.go
	// as that reverses this for log streaming
	switch runtime.EnvCloud {
	case cloud.GCP, cloud.Encore:
		zerolog.LevelFieldName = "severity"
		zerolog.TimestampFieldName = "timestamp"
		zerolog.TimeFieldFormat = time.RFC3339Nano
	default:
	}
}
