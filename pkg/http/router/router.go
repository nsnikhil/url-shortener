package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

func NewRouter(lgr *zap.Logger, newRelic *newrelic.Application, statsdClient reporters.StatsDClient, shortener shortener.Shortener, elongator elongator.Elongator) http.Handler {
	return getChiRouter(lgr, newRelic, statsdClient, shortener, elongator)
}

func getChiRouter(lgr *zap.Logger, newRelic *newrelic.Application, statsdClient reporters.StatsDClient, shortener shortener.Shortener, elongator elongator.Elongator) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(nrgorilla.Middleware(newRelic))

	r.Get(pingPath, handler.PingHandler(lgr, statsdClient))
	r.Post(shortenPath, handler.ShortenHandler(lgr, statsdClient, shortener))
	r.Get(redirectPath, handler.RedirectHandler(lgr, statsdClient, elongator))

	return r
}
