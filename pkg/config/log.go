package config

type LogConfig struct {
	level                  string
	filePath               string
	fileMaxSizeInMb        int
	fileMaxBackups         int
	fileMaxAge             int
	fileWithLocalTimeStamp bool
}

func (lc LogConfig) GetLevel() string {
	return lc.level
}

func (lc LogConfig) GetFilePath() string {
	return lc.filePath
}

func (lc LogConfig) GetFileMaxSizeInMb() int {
	return lc.fileMaxSizeInMb
}

func (lc LogConfig) GetFileMaxBackups() int {
	return lc.fileMaxBackups
}

func (lc LogConfig) GetFileMaxAge() int {
	return lc.fileMaxAge
}

func (lc LogConfig) GetFileWithLocalTimeStamp() bool {
	return lc.fileWithLocalTimeStamp
}

func newLogConfig() LogConfig {
	return LogConfig{
		level:                  getString("LOG_LEVEL"),
		filePath:               getString("LOG_FILE_PATH"),
		fileMaxSizeInMb:        getInt("LOG_FILE_MAX_SIZE_IN_MB"),
		fileMaxBackups:         getInt("LOG_FILE_MAX_BACKUPS"),
		fileMaxAge:             getInt("LOG_FILE_MAX_AGE"),
		fileWithLocalTimeStamp: getBool("LOG_FILE_WITH_LOCAL_TIME_STAMP"),
	}
}
