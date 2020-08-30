package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	env             string
	migrationPath   string
	serverConfig    ServerConfig
	newRelicConfig  NewRelicConfig
	statsDConfig    StatsDConfig
	shortenerConfig ShortenerConfig
	redisConfig     RedisConfig
	databaseConfig  DatabaseConfig
}

func (c Config) GetServerConfig() ServerConfig {
	return c.serverConfig
}

func (c Config) GetNewRelicConfig() NewRelicConfig {
	return c.newRelicConfig
}

func (c Config) GetStatsDConfig() StatsDConfig {
	return c.statsDConfig
}

func (c Config) GetEnv() string {
	return c.env
}

func (c Config) GetMigrationPath() string {
	return c.migrationPath
}

func (c Config) GetShortenerConfig() ShortenerConfig {
	return c.shortenerConfig
}

func (c Config) GetRedisConfig() RedisConfig {
	return c.redisConfig
}

func (c Config) GetDatabaseConfig() DatabaseConfig {
	return c.databaseConfig
}

func NewConfig() Config {
	viper.AutomaticEnv()
	viper.SetConfigName("local")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")
	viper.AddConfigPath("./../../../")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return Config{
		env:             getString("ENV"),
		migrationPath:   getString("MIGRATION_PATH"),
		serverConfig:    newServerConfig(),
		newRelicConfig:  newNewRelicConfig(),
		statsDConfig:    newStatsDConfig(),
		shortenerConfig: newShortenerConfig(),
		redisConfig:     newRedisConfig(),
		databaseConfig:  newDatabaseConfig(),
	}
}
