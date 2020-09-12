package reporters

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"urlshortner/pkg/config"
)

func getLogger(cfg config.LogConfig) *zap.Logger {
	core := zapcore.NewCore(
		getEncoder(),
		getWriteSyncer(cfg),
		getLogLevel(),
	)

	return zap.New(core)
}

func getWriteSyncer(cfg config.LogConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.GetFilePath(),
		MaxSize:    cfg.GetFileMaxSizeInMb(),
		MaxBackups: cfg.GetFileMaxBackups(),
		MaxAge:     cfg.GetFileMaxAge(),
		LocalTime:  cfg.GetFileWithLocalTimeStamp(),
	})
}

func getLogLevel() zapcore.LevelEnabler {
	// TODO PICK FROM CONFIG
	return zap.InfoLevel
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(getLogConfig())
}

func getLogConfig() zapcore.EncoderConfig {
	return zap.NewProductionEncoderConfig()
}
