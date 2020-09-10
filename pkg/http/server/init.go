package server

import (
	"go.uber.org/zap"
	"golang.org/x/crypto/sha3"
	"urlshortner/cmd/config"
	"urlshortner/pkg/elongator"
	"urlshortner/pkg/shortener"
	"urlshortner/pkg/store"
)

type services struct {
	shortener shortener.Shortener
	elongator elongator.Elongator
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
