package initialize

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"lanshan_chat/app/api/global"
	"os"
	"time"
)

func SetupLogger() {
	level := zap.NewAtomicLevel()
	level.SetLevel(zap.DebugLevel)
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:       "msg",
		LevelKey:         "level",
		TimeKey:          "time",
		NameKey:          "logger",
		CallerKey:        "caller",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       CustomTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.FullCallerEncoder,
		ConsoleSeparator: "",
	})

	cores := [...]zapcore.Core{
		zapcore.NewCore(encoder, os.Stdout, level),
		zapcore.NewCore(
			encoder,
			zapcore.AddSync(getwritesync()),
			level,
		),
	}
	global.Logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	defer func(Logger *zap.Logger) {
		_ = global.Logger.Sync()
	}(global.Logger)

	global.Logger.Info("log initialization success")
}

func getwritesync() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   global.Config.ZapConfig.Filename,
		MaxSize:    global.Config.ZapConfig.MaxSize,
		MaxBackups: global.Config.ZapConfig.MaxBackups,
		MaxAge:     global.Config.ZapConfig.MaxAge,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
