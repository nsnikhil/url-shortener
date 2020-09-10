package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

const (
	zero       = 0
	empty      = ""
	insertURL  = "insert into shortener (url, urlhash) values ($1, $2) on conflict on constraint shortener_url_key do nothing"
	getLongURL = "select url from shortener where urlhash = $1"
	getURLHash = "select urlhash from shortener where url = $1"
)

type ShortnerStore interface {
	Save(url, urlHash string) (int, error)
	GetURL(urlHash string) (string, error)
	GetURLHash(url string) (string, error)
	Delete(url, urlHash string) (int, error)
}

type urlShortnerStore struct {
	redis    *redis.Client
	db       *sql.DB
	lgr      *zap.Logger
	newRelic *newrelic.Application
}

func (uss *urlShortnerStore) Save(url, urlHash string) (int, error) {
	count, err := execQuery(uss.lgr, uss.newRelic, "save", uss.db, insertURL, url, urlHash)
	if err != nil {
		return zero, err
	}

	go func(lgr *zap.Logger, redis *redis.Client, url, urlHash string) {
		if err := saveToCache(lgr, redis, url, urlHash); err != nil {
			lgr.Error(err.Error())
		}
	}(uss.lgr, uss.redis, url, urlHash)

	return count, nil
}

func (uss *urlShortnerStore) GetURL(urlHash string) (string, error) {
	return get(getLongURL, urlHash, uss.redis, uss.db, uss.lgr, uss.newRelic, "getURL")
}

func (uss *urlShortnerStore) GetURLHash(url string) (string, error) {
	return get(getURLHash, url, uss.redis, uss.db, uss.lgr, uss.newRelic, "getURLHash")
}

func (uss *urlShortnerStore) Delete(longURL, shortURL string) (int, error) {
	return 0, fmt.Errorf("unimplemented")
}

func get(query, url string, redis *redis.Client, db *sql.DB, lgr *zap.Logger, newRelic *newrelic.Application, name string) (string, error) {
	var res string
	var err error

	res, err = fetchFromCache(redis, url, lgr)
	if err == nil {
		return res, nil
	}

	txn := newRelic.StartTransaction(name)
	ctx := newrelic.NewContext(context.Background(), txn)

	err = db.QueryRowContext(ctx, query, url).Scan(&res)
	if err != nil {
		lgr.Error(err.Error())
		return empty, err
	}

	return res, nil
}

func execQuery(lgr *zap.Logger, newRelic *newrelic.Application, name string, db *sql.DB, query string, args ...interface{}) (int, error) {
	txn := newRelic.StartTransaction(name)
	ctx := newrelic.NewContext(context.Background(), txn)

	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		lgr.Error(err.Error())
		return zero, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		lgr.Error(err.Error())
		return zero, err
	}

	return int(count), nil
}

func saveToCache(lgr *zap.Logger, redis *redis.Client, longURL, shortURL string) error {
	pl := redis.TxPipeline()
	ctx := context.Background()

	//TODO EDIT TTL
	pl.Set(ctx, longURL, shortURL, -1)
	pl.Set(ctx, shortURL, longURL, -1)

	_, err := pl.Exec(ctx)
	if err != nil {
		lgr.Error(err.Error())
		return err
	}

	return nil
}

func fetchFromCache(redis *redis.Client, url string, lgr *zap.Logger) (string, error) {
	ctx := context.Background()

	cmd := redis.Get(ctx, url)

	res, err := cmd.Result()
	if err != nil {
		lgr.Error(err.Error())
		return empty, err
	}

	return res, nil
}

func NewShortnerStore(redis *redis.Client, db *sql.DB, lgr *zap.Logger) ShortnerStore {
	return &urlShortnerStore{
		redis: redis,
		db:    db,
		lgr:   lgr,
	}
}
