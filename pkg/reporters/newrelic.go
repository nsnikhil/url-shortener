package reporters

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"os"
	"urlshortner/pkg/config"
)

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