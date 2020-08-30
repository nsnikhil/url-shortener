package server

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"urlshortner/cmd/config"
	"urlshortner/pkg/shortener/service"
	"urlshortner/pkg/shortener/store"
)

func initService(cfg config.Config, lgr *zap.Logger) *service.Service {
	hashGenerator := service.NewHashGenerator(sha3.New512(), cfg.GetShortenerConfig().GetHashLength())
	urlBuilder := service.NewURLBuilder(cfg.GetShortenerConfig().GetBaseURL())

	ss := service.NewShortenerService(lgr, hashGenerator, urlBuilder, initStore(cfg, lgr))

	return service.NewService(ss)
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
