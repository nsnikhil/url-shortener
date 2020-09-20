package reporters

import (
	"fmt"
	"gopkg.in/alexcesaro/statsd.v2"
	"urlshortner/pkg/config"
)

const (
	attempt = "attempt"
	success = "success"
	failure = "failure"
)

type StatsDClient interface {
	ReportAttempt(bucket string)
	ReportSuccess(bucket string)
	ReportFailure(bucket string)
}

type defaultStatsDClient struct {
	client *statsd.Client
}

func (dsc *defaultStatsDClient) ReportAttempt(bucket string) {
	incBucket(bucket, attempt, dsc.client)
}

func (dsc *defaultStatsDClient) ReportSuccess(bucket string) {
	incBucket(bucket, success, dsc.client)
}

func (dsc *defaultStatsDClient) ReportFailure(bucket string) {
	incBucket(bucket, failure, dsc.client)
}

func (dsc *defaultStatsDClient) NewTiming() statsd.Timing {
	return dsc.client.NewTiming()
}

func incBucket(api, call string, cl *statsd.Client) {
	if cl == nil {
		// TODO LOG THIS
		fmt.Printf("failed to report %s.%s\n", api, call)
		return
	}

	cl.Increment(fmt.Sprintf("%s.%s.counter", api, call))
}

func NewStatsDClient(cfg config.StatsDConfig) StatsDClient {
	return &defaultStatsDClient{
		client: getStatsD(cfg),
	}
}

func getStatsD(cfg config.StatsDConfig) *statsd.Client {
	var err error
	var sc *statsd.Client

	sc, err = statsd.New(statsd.Address(cfg.GetAddress()), statsd.Prefix(cfg.GetNamespace()))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return sc
}
