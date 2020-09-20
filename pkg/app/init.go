package app

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"io"
	"net/http"
	"os"
	"urlshortner/pkg/config"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/http/router"
	"urlshortner/pkg/reporters"
	"urlshortner/pkg/shortener"
	"urlshortner/pkg/store"
)

type services struct {
	shortener shortener.Shortener
	elongator elongator.Elongator
}

func initRouter(cfg config.Config, lgr *zap.Logger, newRelic *newrelic.Application, statsdClient reporters.StatsDClient) http.Handler {
	str := initStore(cfg, lgr, newRelic)
	svc := initService(cfg, lgr, str)

	return router.NewRouter(
		lgr,
		newRelic,
		statsdClient,
		svc.shortener,
		svc.elongator,
	)
}

func initService(cfg config.Config, lgr *zap.Logger, str store.ShortenerStore) *services {
	hashGenerator := shortener.NewHashGenerator(sha3.New512(), cfg.GetShortenerConfig().GetHashLength())
	urlBuilder := shortener.NewURLBuilder(cfg.GetShortenerConfig().GetBaseURL())

	sh := shortener.NewShortener(lgr, str, urlBuilder, hashGenerator)
	el := elongator.NewElongator(lgr, str)

	return &services{shortener: sh, elongator: el}
}

func initStore(cfg config.Config, lgr *zap.Logger, newRelic *newrelic.Application) store.ShortenerStore {
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

	sdb := store.NewShortenerDatabase(db, lgr, newRelic)
	scc := store.NewShortenerCache(cache, lgr)

	return store.NewShortenerStore(sdb, scc, lgr)
}

func initLogger(cfg config.Config) *zap.Logger {
	return reporters.NewLogger(
		cfg.GetEnv(),
		cfg.GetLogConfig().GetLevel(),
		getWriters()...,
	)
}

func getWriters() []io.Writer {
	// TODO ADD LUMBERJACK
	return []io.Writer{
		os.Stdout,
	}
}
