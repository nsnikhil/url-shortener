package main

import (
	"urlshortner/cmd/config"
	server2 "urlshortner/pkg/http/server"
)

func serve() {
	cfg := config.NewConfig()
	rp := newReporters(cfg)
	server2.NewServer(cfg, rp.getLogger(), rp.getNewrelic(), rp.getStatsD()).Start()
}
