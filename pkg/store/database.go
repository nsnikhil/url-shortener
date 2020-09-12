package store

import (
	"context"
	"database/sql"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

type ShortenerDatabase interface {
	Save(url, urlHash string) error
	Get(urlHash string) (string, error)
}

const (
	saveQuery   = "insert into shortener (url, urlhash) values ($1, $2) on conflict on constraint shortener_url_key do nothing returning id"
	getURLQuery = "select url from shortener where urlhash = $1"
)

type urlShortenerDatabase struct {
	lgr      *zap.Logger
	db       *sql.DB
	newRelic *newrelic.Application
}

func (usd *urlShortenerDatabase) Save(url, urlHash string) error {
	return execQuery(usd.lgr, usd.newRelic, "save", usd.db, saveQuery, url, urlHash)
}

func (usd *urlShortenerDatabase) Get(urlHash string) (string, error) {
	return get(getURLQuery, urlHash, usd.db, usd.lgr, usd.newRelic, "getURL")
}

func execQuery(lgr *zap.Logger, newRelic *newrelic.Application, name string, db *sql.DB, query string, args ...interface{}) error {
	txn := newRelic.StartTransaction(name)
	ctx := newrelic.NewContext(context.Background(), txn)

	_, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		lgr.Error(err.Error())
		return err
	}

	return nil
}

func get(query, url string, db *sql.DB, lgr *zap.Logger, newRelic *newrelic.Application, name string) (string, error) {
	var res string
	var err error

	txn := newRelic.StartTransaction(name)
	ctx := newrelic.NewContext(context.Background(), txn)

	err = db.QueryRowContext(ctx, query, url).Scan(&res)
	if err != nil {
		lgr.Error(err.Error())
		return "", err
	}

	return res, nil
}

func NewShortenerDatabase(db *sql.DB, lgr *zap.Logger, newRelic *newrelic.Application) ShortenerDatabase {
	return &urlShortenerDatabase{
		db:       db,
		lgr:      lgr,
		newRelic: newRelic,
	}
}
