package rollinglog

import (
	kZapLogger "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type KZapLogger2 struct {
	zapLogger *zap.Logger
	kLogger   log.Logger
}

func New(opts *Options, kv ...interface{}) *KZapLogger2 {
	if opts == nil {
		opts = NewOptions()
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	zc := zapcore.NewCore(
		generateEncoder(opts),
		multiWriteSyncer(generateWriterSyncer("OutputPaths", opts)...),
		zap.NewAtomicLevelAt(zapLevel),
	)
	l := zap.New(zc, withZapOptions(opts)...)

	logger := &KZapLogger2{
		zapLogger: l.Named(opts.Name),
	}

	zap.RedirectStdLog(l)

	kv = append(kv, "trace_id", tracing.TraceID(), "span_id", tracing.SpanID())
	kLogger := log.With(kZapLogger.NewLogger(logger.zapLogger), kv...)
	logger.kLogger = kLogger

	return logger
}

func (l *KZapLogger2) Sync() {
	_ = l.zapLogger.Sync()
}

func (l KZapLogger2) GetLogger() log.Logger {
	return l.kLogger
}

func generateEncoder(opts *Options) (encode zapcore.Encoder) {
	encodeLevel := zapcore.CapitalLevelEncoder
	// when output to local path and format = "console", with color is forbidden.
	if opts.Format == consoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		FunctionKey:    "func",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	switch opts.Format {
	case "console":
		encode = zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		encode = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encode = zapcore.NewJSONEncoder(encoderConfig)
	}
	return
}

func withZapOptions(opts *Options) []zap.Option {
	var zOpts []zap.Option
	zOpts = append(zOpts, zap.WithCaller(opts.DisableCaller))
	zOpts = append(zOpts, zap.AddStacktrace(zapcore.ErrorLevel))
	zOpts = append(zOpts, zap.AddCallerSkip(4))
	zOpts = append(zOpts, zap.ErrorOutput(multiWriteSyncer(generateWriterSyncer("ErrorOutputPaths", opts)...)))

	if opts.Development {
		zOpts = append(zOpts, zap.Development())
	}

	return zOpts
}
