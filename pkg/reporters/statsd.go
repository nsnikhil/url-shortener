package reporters

import (
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"urlshortner/pkg/config"
)

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
