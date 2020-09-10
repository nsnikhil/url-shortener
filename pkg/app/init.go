package app

import (
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"gopkg.in/alexcesaro/statsd.v2"
	"urlshortner/pkg/config"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/http/router"
	"urlshortner/pkg/shortener"
	"urlshortner/pkg/store"
)

type services struct {
	shortener shortener.Shortener
	elongator elongator.Elongator
}

func initRouter(cfg config.Config, lgr *zap.Logger, newRelic *newrelic.Application, statsd *statsd.Client) *mux.Router {
	svc := initService(cfg, lgr)
	return router.NewRouter(
		lgr,
		newRelic,
		statsd,
		svc.shortener,
		svc.elongator,
	)
}

func initService(cfg config.Config, lgr *zap.Logger) *services {
	str := initStore(cfg, lgr)

	hashGenerator := shortener.NewHashGenerator(sha3.New512(), cfg.GetShortenerConfig().GetHashLength())
	urlBuilder := shortener.NewURLBuilder(cfg.GetShortenerConfig().GetBaseURL())

	sh := shortener.NewShortener(lgr, str, urlBuilder, hashGenerator)
	el := elongator.NewElongator(lgr, str)

	return &services{shortener: sh, elongator: el}
}

func initStore(cfg config.Config, lgr *zap.Logger) *store.Store {
	cacheHandler := store.NewCacheHandler(cfg.GetRedisConfig(), lgr)
	cache, err := cacheHandler.GetCache()
	if err != nil {
		lgr.Fatal(err.Error())
	}

	dbHandler := store.NewDBHandler(cfg.GetDatabaseConfig(), lgr)
	db, err := dbHandler.GetDB()
	if err != nil {
		lgr.Fatal(err.Error())
	}

	return store.NewStore(store.NewShortnerStore(cache, db, lgr))
}
