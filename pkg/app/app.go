package app

import (
	"urlshortner/pkg/config"
	"urlshortner/pkg/http/server"
	"urlshortner/pkg/reporters"
)

func Start() {
	cfg := config.NewConfig()
	rp := reporters.NewReporters(cfg)
	rt := initRouter(cfg, rp.GetLogger(), rp.GetNewrelic(), rp.GetStatsD())
	server.NewServer(cfg, rp.GetLogger(), rt).Start()
}
