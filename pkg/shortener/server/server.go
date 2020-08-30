package server

import (
	"context"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlshortner/cmd/config"
	"urlshortner/pkg/shortener/router"
)

type Server interface {
	Start()
}

type appServer struct {
	cfg      config.Config
	lgr      *zap.Logger
	newRelic *newrelic.Application
	statsd   *statsd.Client
}

func (as *appServer) Start() {
	svc := initService(as.cfg, as.lgr)
	rt := router.NewRouter(as.lgr, as.newRelic, as.statsd, svc)
	server := newHTTPServer(as.cfg.GetServerConfig(), rt)

	as.lgr.Sugar().Infof("listening on %s", as.cfg.GetServerConfig().GetAddress())
	go func() { _ = server.ListenAndServe() }()

	waitForShutdown(server, as.lgr)
}

func waitForShutdown(server *http.Server, lgr *zap.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigCh

	err := server.Shutdown(context.Background())
	if err != nil {
		lgr.Error(err.Error())
		return
	}

	lgr.Info("server shutdown successful")
}

func newHTTPServer(cfg config.ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Handler:      handler,
		Addr:         cfg.GetAddress(),
		WriteTimeout: time.Second * time.Duration(cfg.GetReadTimeout()),
		ReadTimeout:  time.Second * time.Duration(cfg.GetWriteTimeout()),
	}
}

func NewServer(cfg config.Config, lgr *zap.Logger, newRelic *newrelic.Application, statsd *statsd.Client) Server {
	return &appServer{
		cfg:      cfg,
		lgr:      lgr,
		newRelic: newRelic,
		statsd:   statsd,
	}
}
