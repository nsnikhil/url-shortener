package main

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"os"
	"urlshortner/cmd/config"
)

const dev = "dev"

type reporters struct {
	logger *zap.Logger
	nrApp  *newrelic.Application
	sc     *statsd.Client
}

func (rp *reporters) getLogger() *zap.Logger {
	return rp.logger
}

func (rp *reporters) getNewrelic() *newrelic.Application {
	return rp.nrApp
}

func (rp *reporters) getStatsD() *statsd.Client {
	return rp.sc
}

func newReporters(cfg config.Config) *reporters {
	return &reporters{
		logger: getLogger(cfg.GetEnv()),
		nrApp:  getNewRelic(cfg.GetNewRelicConfig()),
		sc:     getStatsD(cfg.GetStatsDConfig()),
	}
}

func getLogger(env string) *zap.Logger {
	var err error
	var lgr *zap.Logger

	if env == dev {
		lgr, err = zap.NewDevelopment()
	} else {
		lgr, err = zap.NewProduction()
	}

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer func() {
		if err := lgr.Sync(); err != nil {
			fmt.Println(err)
		}
	}()

	return lgr
}

func getNewRelic(nrc config.NewRelicConfig) *newrelic.Application {
	var err error
	var nrApp *newrelic.Application

	nrApp, err = newrelic.NewApplication(
		newrelic.ConfigAppName(nrc.GetAppName()),
		newrelic.ConfigLicense(nrc.GetLicenseKey()),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nrApp
}

func getStatsD(sdc config.StatsDConfig) *statsd.Client {
	var err error
	var sc *statsd.Client

	sc, err = statsd.New(statsd.Address(sdc.GetAddress()), statsd.Prefix(sdc.GetNamespace()))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return sc
}
