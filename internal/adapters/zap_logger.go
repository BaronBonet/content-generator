package adapters

import (
	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/content-generator/internal/infrastructure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(config zap.Config, useDebug bool) ports.Logger {

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	log, err := config.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic("Cannot create logger: " + err.Error())
	}
	log = log.WithOptions(zap.AddCallerSkip(2))

	if useDebug {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	return &zapLogger{
		logger: log.With(zap.String("version", infrastructure.Version)),
	}
}

func (l *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log(zap.DebugLevel, msg, keysAndValues...)
}

func (l *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log(zap.InfoLevel, msg, keysAndValues...)
}

func (l *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log(zap.WarnLevel, msg, keysAndValues...)
}

func (l *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log(zap.ErrorLevel, msg, keysAndValues...)
}

func (l *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log(zap.ErrorLevel, msg, keysAndValues...)
	panic(msg)
}

func (l *zapLogger) log(level zapcore.Level, msg string, keysAndValues ...interface{}) {
	if l.logger.Check(level, msg) == nil {
		return
	}
	fields := createFields(keysAndValues)
	l.logger.Log(level, msg, fields...)
}

// createFields creates zap fields from a list of keys and values.
func createFields(values []interface{}) []zap.Field {
	fields := make([]zap.Field, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		fields[i/2] = zap.Any(values[i].(string), values[i+1])
	}
	return fields
}
