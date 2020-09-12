package store

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

type ShortenerCache interface {
	Save(url, urlHash string, ttl int) error
	Get(urlHash string) (string, error)
}

type urlShortenerCache struct {
	client *redis.Client
	lgr    *zap.Logger
}

func (usc *urlShortenerCache) Save(url, urlHash string, ttl int) error {
	ctx := context.Background()

	cmd := usc.client.Set(ctx, urlHash, url, time.Second*time.Duration(ttl))
	if cmd.Err() != nil {
		usc.lgr.Error(cmd.Err().Error())
		return cmd.Err()
	}

	return nil
}

func (usc *urlShortenerCache) Get(urlHash string) (string, error) {
	ctx := context.Background()

	cmd := usc.client.Get(ctx, urlHash)

	res, err := cmd.Result()
	if err != nil {
		usc.lgr.Error(err.Error())
		return "", err
	}

	return res, nil
}

func NewShortenerCache(client *redis.Client, lgr *zap.Logger) ShortenerCache {
	return &urlShortenerCache{
		client: client,
		lgr:    lgr,
	}
}
