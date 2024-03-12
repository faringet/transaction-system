package zaplogger

import (
	"errors"
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger Modes
const (
	Production  Mode = "production"
	Development Mode = "development"
	None        Mode = "none"
)

var ErrUnsupportedZapLoggerMode = errors.New("unsupported zapLogger mode")

type Mode string

// Field Keys
const (
	payloadKey = "payload"
)

func New(mode Mode) (logger *zap.Logger, cleanup func(), err error) {
	var zapLogger *zap.Logger

	switch mode {
	case Development:
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, err = config.Build()
	case Production:
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLogger, err = config.Build()
	case None:
		zapLogger = zap.NewNop()
	default:
		err = fmt.Errorf("%w: %s", ErrUnsupportedZapLoggerMode, mode)
	}

	if err != nil {
		return nil, nil, err
	}

	undoRedirectStdLog := zap.RedirectStdLog(zapLogger)
	cleanup = func() {
		if errSync := zapLogger.Sync(); errSync != nil {
			log.Println(errSync)
		}

		undoRedirectStdLog()
	}

	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(0)).With(zap.Namespace(payloadKey))

	return zapLogger, cleanup, nil
}
