package main

import (
	"urlshortner/cmd/config"
	"urlshortner/pkg/shortener/server"
)

func serve() {
	cfg := config.NewConfig()
	rp := newReporters(cfg)
	server.NewServer(cfg, rp.getLogger(), rp.getNewrelic(), rp.getStatsD()).Start()
}
