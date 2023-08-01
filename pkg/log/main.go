package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(s *Logger)

func WithLevel(level zapcore.Level) Option {
	return func(l *Logger) {
		l.atom.SetLevel(level)

		if level == zap.DebugLevel {
			l.Logger = *l.Logger.WithOptions(zap.AddCaller())
		}
	}
}

func WithGlobalLogger() Option {
	return func(l *Logger) {
		zap.ReplaceGlobals(&l.Logger)
	}
}

type Logger struct {
	zap.Logger
	atom zap.AtomicLevel
}

func NewLogger(options ...Option) (*Logger, error) {
	atom := zap.NewAtomicLevelAt(zap.WarnLevel)
	econfig := zap.NewProductionEncoderConfig()
	econfig.EncodeTime = func(ts time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(ts.UTC().Format(time.RFC3339))
	}
	econfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config := zap.Config{
		Level:            atom,
		Encoding:         "console",
		EncoderConfig:    econfig,
		DisableCaller:    true,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{},
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	l := &Logger{*logger, atom}

	for _, option := range options {
		option(l)
	}

	return l, nil
}
