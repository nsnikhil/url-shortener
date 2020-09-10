package config

import (
	"fmt"
)

type StatsDConfig struct {
	host      string
	port      string
	namespace string
}

func newStatsDConfig() StatsDConfig {
	return StatsDConfig{
		host:      getString("STATSD_HOST"),
		port:      getString("STATSD_PORT"),
		namespace: getString("STATSD_NAMESPACE"),
	}
}

func (sdc StatsDConfig) GetNamespace() string {
	return sdc.namespace
}

func (sdc StatsDConfig) GetAddress() string {
	return fmt.Sprintf("%s:%s", sdc.host, sdc.port)
}
