package app

import (
	"urlshortner/pkg/config"
	"urlshortner/pkg/http/server"
	"urlshortner/pkg/reporters"
)

func Start() {
	cfg := config.NewConfig()
	lgr := initLogger(cfg)
	nr := reporters.NewNewRelicApp(cfg.GetNewRelicConfig())
	sd := reporters.NewStatsDClient(cfg.GetStatsDConfig())

	rt := initRouter(cfg, lgr, nr, sd)
	server.NewServer(cfg, lgr, rt).Start()
}
