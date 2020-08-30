package config

import (
	"fmt"
)

type ServerConfig struct {
	host             string
	port             string
	readTimoutInSec  int
	writeTimoutInSec int
}

func newServerConfig() ServerConfig {
	return ServerConfig{
		host:             getString("APP_HOST"),
		port:             getString("APP_PORT"),
		readTimoutInSec:  getInt("APP_READ_TIMEOUT_IN_SEC"),
		writeTimoutInSec: getInt("APP_WRITE_TIMEOUT_IN_SEC"),
	}
}

func (sc ServerConfig) GetAddress() string {
	return fmt.Sprintf(":%s", sc.port)
}

func (sc ServerConfig) GetReadTimeout() int {
	return sc.readTimoutInSec
}

func (sc ServerConfig) GetWriteTimeout() int {
	return sc.readTimoutInSec
}
