package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var lvl int

func SetLogLevel(l int) {
	lvl = l
}

func GetLogLevel() zap.AtomicLevel {
	switch lvl {
	case 1:
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case 2:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case 3:
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

func GetLogger(role string) *zap.Logger {
	cfg := zap.Config{
		Level:            GetLogLevel(),
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"role": role,
		},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			TimeKey:     "ts",
			LevelKey:    "lvl",
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
		},
	}

	logger, _ := cfg.Build()

	return logger
}
