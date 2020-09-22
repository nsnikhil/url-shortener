package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/newrelic/go-agent/v3/integrations/nrgorilla"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/http/internal/handler"
	mdl "urlshortner/pkg/http/internal/middleware"
	"urlshortner/pkg/reporters"
	"urlshortner/pkg/shortener"
)

const (
	pingAPI     = "ping"
	shortenAPI  = "shorten"
	redirectAPI = "redirect"

	pingPath     = "/ping"
	shortenPath  = "/shorten"
	redirectPath = "/{hash_code:[a-zA-Z]+}"
	metricPath   = "/metrics"
)

func NewRouter(lgr *zap.Logger, newRelic *newrelic.Application, prometheus reporters.Prometheus, shortener shortener.Shortener, elongator elongator.Elongator) http.Handler {
	return getChiRouter(lgr, newRelic, prometheus, shortener, elongator)
}

func getChiRouter(lgr *zap.Logger, newRelic *newrelic.Application, pr reporters.Prometheus, shortener shortener.Shortener, elongator elongator.Elongator) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(nrgorilla.Middleware(newRelic))

	sh := handler.NewShortenHandler(shortener)
	rh := handler.NewRedirectHandler(elongator)

	r.Get(pingPath, withMiddlewares(lgr, pr, pingAPI, handler.PingHandler()))
	r.Post(shortenPath, withMiddlewares(lgr, pr, shortenAPI, mdl.WithError(sh.Shorten)))
	r.Handle(metricPath, promhttp.Handler())

	r.Get(redirectPath, withMiddlewares(lgr, pr, redirectAPI, mdl.WithError(rh.Redirect)))

	return r
}

func withMiddlewares(lgr *zap.Logger, prometheus reporters.Prometheus, api string, handler func(resp http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return mdl.WithReqRespLog(lgr,
		mdl.WithResponseHeaders(
			mdl.WithPrometheus(prometheus, api, handler),
		),
	)
}
