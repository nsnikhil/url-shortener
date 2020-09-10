package reporters

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"urlshortner/pkg/config"
)

type Reporters struct {
	logger *zap.Logger
	nrApp  *newrelic.Application
	sc     *statsd.Client
}

func (rp *Reporters) GetLogger() *zap.Logger {
	return rp.logger
}

func (rp *Reporters) GetNewrelic() *newrelic.Application {
	return rp.nrApp
}

func (rp *Reporters) GetStatsD() *statsd.Client {
	return rp.sc
}

func NewReporters(cfg config.Config) *Reporters {
	return &Reporters{
		logger: getLogger(cfg.GetEnv()),
		nrApp:  getNewRelic(cfg.GetNewRelicConfig()),
		sc:     getStatsD(cfg.GetStatsDConfig()),
	}
}
