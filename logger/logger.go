package logger

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewgomieDevelopmentEncoderConfig returns an opinionated EncoderConfig for
// development environments.
func NewgomieDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// NewDevelopmentConfig is a reasonable development logging configuration.
// Logging is enabled at DebugLevel and above.
//
// It enables development mode (which makes DPanicLevel logs panic), uses a
// console encoder, writes to standard error, and disables sampling.
// Stacktraces are automatically included on logs of WarnLevel and above.
func NewDevelopmentConfig() zap.Config {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    NewgomieDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

type Logger = zap.SugaredLogger

type LogCtx = struct {
	level        string
	message      string
	packageName  string
	functionName string
	traceIid     string
}

var logger *Logger

var theLogger *Logger

func GetLogger() *Logger {
	/* The zap logger is fast and the sugared variant makes simple JSON as key-value pairs
	 */
	if theLogger == nil {
		// cfg := config.GetConfig()

		logger, _ := zap.NewDevelopment() // or NewProduction or NewDevelopment or NewExample
		defer logger.Sync()
		theLogger = logger.Sugar()
	}
	return theLogger
}

func LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if theLogger == nil {
			GetLogger()
		}
		theLogger.Infof(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
