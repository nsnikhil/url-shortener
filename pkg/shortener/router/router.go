package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"gopkg.in/alexcesaro/statsd.v2"
	"net/http"
	"urlshortner/pkg/shortener/handler"
	"urlshortner/pkg/shortener/service"
)

const (
	pingPath     = "/ping"
	shortenPath  = "/shorten"
	redirectPath = "/{hash_code}"
)

func NewRouter(lgr *zap.Logger, newRelic *newrelic.Application, statsd *statsd.Client, svc *service.Service) *mux.Router {
	r := mux.NewRouter()
	r.Use(nrgorilla.Middleware(newRelic))
	r.Use(handlers.RecoveryHandler())

	r.HandleFunc(pingPath, handler.PingHandler(lgr, statsd)).Methods(http.MethodGet)
	r.HandleFunc(shortenPath, handler.ShortenHandler(lgr, statsd, svc)).Methods(http.MethodPost)
	r.HandleFunc(redirectPath, handler.RedirectHandler(lgr, statsd, svc)).Methods(http.MethodGet)

	return r
}
