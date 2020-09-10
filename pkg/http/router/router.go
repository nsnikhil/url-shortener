package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"net/http"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/http/handler"
	"urlshortner/pkg/reporters"
	"urlshortner/pkg/shortener"
)

const (
	pingPath     = "/ping"
	shortenPath  = "/shorten"
	redirectPath = "/{hash_code}"
)

func NewRouter(lgr *zap.Logger, newRelic *newrelic.Application, statsdClient reporters.StatsDClient, shortener shortener.Shortener, elongator elongator.Elongator) *mux.Router {
	r := mux.NewRouter()
	r.Use(nrgorilla.Middleware(newRelic))
	r.Use(handlers.RecoveryHandler())

	r.HandleFunc(pingPath, handler.PingHandler(lgr, statsdClient)).Methods(http.MethodGet)
	r.HandleFunc(shortenPath, handler.ShortenHandler(lgr, statsdClient, shortener)).Methods(http.MethodPost)
	r.HandleFunc(redirectPath, handler.RedirectHandler(lgr, statsdClient, elongator)).Methods(http.MethodGet)

	return r
}
