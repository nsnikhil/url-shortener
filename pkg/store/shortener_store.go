package store

import (
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

// TODO: READ FROM CONFIG
const ttlInSec = 3600

type ShortenerStore interface {
	Save(url, urlHash string) error
	GetURL(urlHash string) (string, error)
}

type urlShortnerStore struct {
	db    ShortenerDatabase
	cache ShortenerCache

	lgr      *zap.Logger
	newRelic *newrelic.Application
}

func (uss *urlShortnerStore) Save(url, urlHash string) error {
	err := uss.db.Save(url, urlHash)
	if err != nil {
		return err
	}

	go func(cache ShortenerCache, lgr *zap.Logger, url, urlHash string, ttl int) {
		if err := cache.Save(url, urlHash, ttl); err != nil {
			lgr.Error(err.Error())
		}
	}(uss.cache, uss.lgr, url, urlHash, ttlInSec)

	return nil
}

func (uss *urlShortnerStore) GetURL(urlHash string) (string, error) {
	res, err := uss.cache.Get(urlHash)
	if err == nil {
		return res, nil
	}

	res, err = uss.db.Get(urlHash)
	if err != nil {
		return "", err
	}

	return res, nil
}

func NewShortenerStore(db ShortenerDatabase, cache ShortenerCache, lgr *zap.Logger) ShortenerStore {
	return &urlShortnerStore{
		db:    db,
		cache: cache,
		lgr:   lgr,
	}
}
