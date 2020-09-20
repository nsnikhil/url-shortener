package config

type LogConfig struct {
	level string
}

func (lc LogConfig) GetLevel() string {
	return lc.level
}

func newLogConfig() LogConfig {
	return LogConfig{
		level: getString("LOG_LEVEL"),
	}
}
