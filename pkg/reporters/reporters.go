package reporters

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"urlshortner/pkg/config"
)

type Reporters struct {
	logger *zap.Logger
	nrApp  *newrelic.Application
	sc     StatsDClient
}

func (rp *Reporters) GetLogger() *zap.Logger {
	return rp.logger
}

func (rp *Reporters) GetNewrelic() *newrelic.Application {
	return rp.nrApp
}

func (rp *Reporters) GetStatsD() StatsDClient {
	return rp.sc
}

func NewReporters(cfg config.Config) *Reporters {
	return &Reporters{
		logger: getLogger(cfg.GetEnv()),
		nrApp:  getNewRelic(cfg.GetNewRelicConfig()),
		sc:     newStatsDClient(getStatsD(cfg.GetStatsDConfig())),
	}
}
