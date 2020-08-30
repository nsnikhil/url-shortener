package config

import "fmt"

type RedisConfig struct {
	host              string
	port              string
	password          string
	db                int
	maxRetry          int
	dialTimeoutInSec  int
	readTimeoutInSec  int
	writeTimeoutInSec int
	poolSize          int
	minIdleConnection int
	keyExpirationTTL  int
}

func newRedisConfig() RedisConfig {
	return RedisConfig{
		host:              getString("REDIS_HOST"),
		port:              getString("REDIS_PORT"),
		password:          getString("REDIS_PASSWORD"),
		db:                getInt("REDIS_DB"),
		maxRetry:          getInt("REDIS_MAX_RETRY"),
		dialTimeoutInSec:  getInt("REDIS_DIAL_TIMOUT_IN_SEC"),
		readTimeoutInSec:  getInt("REDIS_READ_TIMOUT_IN_SEC"),
		writeTimeoutInSec: getInt("REDIS_WRITE_TIMOUT_IN_SEC"),
		poolSize:          getInt("REDIS_POOL_SIZE"),
		minIdleConnection: getInt("REDIS_MIN_IDLE_CONNECTION"),
		keyExpirationTTL:  getInt("REDIS_KEY_EXPIRATION_TIME_IN_HOURS"),
	}
}

func (rc RedisConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", rc.host, rc.port)
}

func (rc RedisConfig) GetPassword() string {
	return rc.password
}

func (rc RedisConfig) GetDB() int {
	return rc.db
}

func (rc RedisConfig) GetMaxRetry() int {
	return rc.maxRetry
}

func (rc RedisConfig) GetDialTimeoutInSec() int {
	return rc.dialTimeoutInSec
}

func (rc RedisConfig) GetReadTimeoutInSec() int {
	return rc.readTimeoutInSec
}

func (rc RedisConfig) GetWriteTimeoutInSec() int {
	return rc.writeTimeoutInSec
}

func (rc RedisConfig) GetPoolSize() int {
	return rc.poolSize
}

func (rc RedisConfig) GetMinIdleConnection() int {
	return rc.minIdleConnection
}

func (rc RedisConfig) GetKeyExpirationTTL() int {
	return rc.keyExpirationTTL
}
