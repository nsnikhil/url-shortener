package reporters

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"urlshortner/pkg/config"
)

// TODO
func logFile(cfg config.LogConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   cfg.GetFilePath(),
		MaxSize:    cfg.GetFileMaxSizeInMb(),
		MaxBackups: cfg.GetFileMaxBackups(),
		MaxAge:     cfg.GetFileMaxAge(),
		LocalTime:  cfg.GetFileWithLocalTimeStamp(),
	}
}
