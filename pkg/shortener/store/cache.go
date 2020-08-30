package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
	"urlshortner/cmd/config"
)

type CacheHandler interface {
	GetCache() (*redis.Client, error)
}

type defaultCacheHandler struct {
	cfg config.RedisConfig
	lgr *zap.Logger
}

func (dh *defaultCacheHandler) GetCache() (*redis.Client, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:         dh.cfg.GetAddress(),
		Password:     dh.cfg.GetPassword(),
		DB:           dh.cfg.GetDB(),
		MaxRetries:   dh.cfg.GetMaxRetry(),
		DialTimeout:  time.Second * time.Duration(dh.cfg.GetDialTimeoutInSec()),
		ReadTimeout:  time.Second * time.Duration(dh.cfg.GetReadTimeoutInSec()),
		WriteTimeout: time.Second * time.Duration(dh.cfg.GetWriteTimeoutInSec()),
		PoolSize:     dh.cfg.GetPoolSize(),
		MinIdleConns: dh.cfg.GetMinIdleConnection(),
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		dh.lgr.Error(err.Error())
		return nil, err
	}

	return rdb, nil
}

func NewCacheHandler(cfg config.RedisConfig, lgr *zap.Logger) CacheHandler {
	return &defaultCacheHandler{
		cfg: cfg,
		lgr: lgr,
	}
}
