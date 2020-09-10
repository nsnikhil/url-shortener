package server

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlshortner/pkg/config"
)

type Server interface {
	Start()
}

type appServer struct {
	cfg    config.Config
	lgr    *zap.Logger
	router http.Handler
}

func (as *appServer) Start() {
	server := newHTTPServer(as.cfg.GetServerConfig(), as.router)

	as.lgr.Sugar().Infof("listening on %s", as.cfg.GetServerConfig().GetAddress())
	go func() { _ = server.ListenAndServe() }()

	waitForShutdown(server, as.lgr)
}

func waitForShutdown(server *http.Server, lgr *zap.Logger) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigCh

	defer func() { _ = lgr.Sync() }()

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

func NewServer(cfg config.Config, lgr *zap.Logger, router http.Handler) Server {
	return &appServer{
		cfg:    cfg,
		lgr:    lgr,
		router: router,
	}
}
