package logtool

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	ServiceName string
	jsonFormat  bool
	sugar       *zap.SugaredLogger
)

func getLogger() *zap.SugaredLogger {
	return sugar
}

func Debug(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

func Init(serviceName string, useJsonFormat bool) {
	ServiceName = serviceName
	jsonFormat = useJsonFormat
	var logger *zap.Logger
	var err error
	if !useJsonFormat {
		zapConfig := zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.TimeKey = ""
		zapConfig.OutputPaths = []string{"stdout"}

		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		zapConfig.EncoderConfig.EncodeName = zapcore.FullNameEncoder
		zapConfig.DisableStacktrace = true

		// Enable caller information and skip logtool package
		zapConfig.Development = true
		zapConfig.EncoderConfig.CallerKey = "caller"
		zapConfig.EncoderConfig.StacktraceKey = "stacktrace"

		logger, err = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		zapConfig := zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "time"
		zapConfig.Encoding = "json"
		zapConfig.OutputPaths = []string{"stdout"}

		// Enable caller information and skip logtool package
		zapConfig.EncoderConfig.CallerKey = "caller"
		zapConfig.EncoderConfig.StacktraceKey = "stacktrace"

		logger, err = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	if err != nil {
		log.Fatal(err)
	}
	sugar = logger.Sugar()
}
