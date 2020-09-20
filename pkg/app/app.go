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
	pr := reporters.NewPrometheus()

	rt := initRouter(cfg, lgr, nr, pr)
	server.NewServer(cfg, lgr, rt).Start()
}
